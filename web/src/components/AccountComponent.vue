<script setup lang="ts">
import { ref, watch, computed, onMounted } from "vue";
import { GatewayError } from "../constants";
import * as EmailValidator from "email-validator";
import {
  User,
  Group,
  GroupList,
  NewUserRequest,
  UpdateUserRequest,
} from "ldap-manager";

import { useRouter } from "vue-router";
import { useToast } from "bootstrap-vue-3";
import { useAuthStore } from "../stores/auth";
import { useGroupsStore } from "../stores/groups";
import { useMembersStore } from "../stores/members";
import { useAppStore } from "../stores/app";
import { useAccountsStore } from "../stores/accounts";

const router = useRouter();
const toast = useToast();
const appStore = useAppStore();
const authStore = useAuthStore();
const accountsStore = useAccountsStore();
const groupsStore = useGroupsStore();
const membersStore = useMembersStore();

const props = withDefaults(
  defineProps<{
    account?: string;
    title: string;
    all: boolean;
    create: boolean;
  }>(),
  {
    title: "Account",
    all: false,
    create: false,
  }
);

const invalidGroupText = ref<string>("no such group");
const error = ref<string | undefined>(undefined);
const groupMemberError = ref<string | undefined>(undefined);
const submissionError = ref<string | undefined>(undefined);

const availableGroups = ref<Group[]>([]);
const userGroups = ref<Group[]>([]);
const userGroupNames = ref<string[]>([]);
const userGroupInputDisabled = ref<boolean>(false);
const groupState = computed(() => true);

const processing = ref<boolean>(false);
const watchGroups = ref<boolean>(false);
const checkingGroup = ref<boolean>(false);

const newUserRequest = ref<NewUserRequest>(NewUserRequest.fromJSON({}));

const passwordConfirm = ref<string>("");

const activeIsAdmin = computed(() => authStore.isAdmin);

watch(userGroupNames, async (after: string[], before: string[]) => {
  if (watchGroups.value) {
    try {
      // lock
      watchGroups.value = false;
      userGroupInputDisabled.value = true;
      processing.value = true;

      const b = new Set(before);
      const a = new Set(after);
      const added = new Set([...a].filter((x) => !b.has(x)));
      const removed = new Set([...b].filter((x) => !a.has(x)));
      console.log("added", added);
      console.log("removed", removed);

      for (let group of added) {
        if (props.account) {
          await addToGroup(props.account, group);
        }
      }
      for (let group of removed) {
        if (props.account) {
          await removeFromGroup(props.account, group);
        }
      }
    } finally {
      watchGroups.value = true;
      userGroupInputDisabled.value = false;
      processing.value = false;
    }
  }
});

function groupValidator(name: string) {
  invalidGroupText.value = "no such group";
  // find candidates
  const matches = availableGroups.value
    .map((group: Group) => group.name)
    .filter((group: string) =>
      group.toLowerCase().includes(name.toLowerCase())
    );
  // find exact match
  const found =
    matches.find(
      (group: string) => group.toLowerCase() == name.toLowerCase()
    ) !== undefined;
  if (!found && matches.length > 0) {
    invalidGroupText.value = `have group ${matches[0]}, but not ${name}`;
  }
  return found;
}

const passwordStrengthVariant = computed(() => {
  if (passwordStrength.value < 3) return "danger";
  if (passwordStrength.value < 6) return "warning";
  return "success";
});

const passwordStrength = computed(() => {
  const pw = newUserRequest.value.password;
  // at least 8 characters
  const sufficientLengthScore = Number(/.{8,}/.test(pw));
  // bonus if longer
  const goodLengthScore = Number(/.{12,}/.test(pw));
  // a lower letter
  const lowercaseScore = Number(/[a-z]/.test(pw));
  // a upper letter
  const uppercaseScore = Number(/[A-Z]/.test(pw));
  // a digit
  const hasDigitScore = Number(/\d/.test(pw));
  // a special character
  const hasSpecialScore = Number(/[^A-Za-z0-9]/.test(pw));
  return (
    1 +
    sufficientLengthScore *
      (goodLengthScore +
        lowercaseScore +
        uppercaseScore +
        hasDigitScore +
        hasSpecialScore)
  );
});

const enteredPassword = computed(() => {
  return (newUserRequest.value.password + passwordConfirm.value).length > 0;
});

const passwordsMatch = computed(() => {
  return newUserRequest.value.password == passwordConfirm.value;
});

const validEmail = computed(() => {
  return EmailValidator.validate(newUserRequest.value.email);
});

const passwordStrengthLabel = computed(() => {
  if (passwordStrength.value < 3) return "weak!";
  if (passwordStrength.value < 6) return "fair enough";
  return "good";
});

function successAlert(message: string) {
  toast?.success(
    {
      title: "Success",
      body: message,
    },
    {
      autoHide: true,
      delay: 5000,
    }
  );
}

async function deleteAccount(username: string) {
  await appStore.newConfirmation({
    message: "Are you sure?",
    ack: "Yes, delete",
  });
  try {
    submissionError.value = undefined;
    processing.value = true;

    await accountsStore.deleteAccount(username);
    successAlert(`${username} was successfully deleted`);
    router.push({ name: "LoginRoute" });
  } catch (err: unknown) {
    if (err instanceof GatewayError) {
      submissionError.value = err.message;
    } else {
      throw err;
    }
  } finally {
    processing.value = false;
  }
}

async function createAccount() {
  if (newUserRequest.value.password !== passwordConfirm.value) return;
  try {
    submissionError.value = undefined;
    processing.value = true;

    await accountsStore.newAccount(newUserRequest.value);
    successAlert(`${newUserRequest.value.username} was created`);
    // edit the user
    router.push({
      name: "EditUserRoute",
      params: { username: newUserRequest.value.username },
    });
  } catch (err: unknown) {
    if (err instanceof GatewayError) {
      submissionError.value = err.message;
    } else {
      throw err;
    }
  } finally {
    processing.value = false;
  }
}

async function removeFromGroup(username: string, group: string) {
  try {
    groupMemberError.value = undefined;

    await membersStore.removeGroupMember({
      username,
      group,
    });
    userGroups.value = userGroups.value.filter((g: Group) => g.name !== group);
    successAlert(`${username} was removed from ${group}`);
  } catch (err: unknown) {
    // add group again
    userGroupNames.value.push(group);
    if (err instanceof GatewayError) {
      groupMemberError.value = err.message;
    } else {
      throw err;
    }
  }
}

async function addToGroup(username: string, group: string) {
  try {
    groupMemberError.value = undefined;

    await membersStore.addGroupMember({
      username: username,
      group: group,
    });

    // add to userGroups
    const addedGroup = availableGroups.value.find(
      (g: Group) => g.name === group
    );
    if (addedGroup) {
      userGroups.value = [...userGroups.value, addedGroup];
    }
    successAlert(`${username} was added to ${group}`);
  } catch (err: unknown) {
    // remove group again
    userGroupNames.value = userGroupNames.value.filter(
      (g: string) => g !== group
    );
    if (err instanceof GatewayError) {
      groupMemberError.value = err.message;
    } else {
      throw err;
    }
  }
}

async function updateAccount(username: string | undefined = undefined) {
  const oldUsername = username ?? props.account;
  if (!oldUsername) {
    return;
  }

  if (newUserRequest.value.password !== passwordConfirm.value) return;

  try {
    processing.value = true;
    submissionError.value = undefined;

    const request: UpdateUserRequest = {
      update: newUserRequest.value,
      username: oldUsername,
    };

    await accountsStore.updateAccount(request);
    successAlert(`${oldUsername} was updated`);

    const updatedUser: User | undefined = await accountsStore.getAccount(
      newUserRequest.value.username ?? oldUsername
    );
    if (!updatedUser) {
      submissionError.value = "invalid user";
      return;
    }
    authStore.updateUser(updatedUser);
  } catch (err: unknown) {
    if (err instanceof GatewayError) {
      submissionError.value = err.message;
    } else {
      throw err;
    }
  } finally {
    processing.value = false;
  }
}

async function onSubmit() {
  props.create ? await createAccount() : await updateAccount();
}

async function fetchAvailableGroups() {
  try {
    error.value = undefined;

    let page = 1;
    let total = null;

    while (total === null || availableGroups.value.length < total) {
      const request = {
        page,
        perPage: 100,
        search: "",
      };
      const list: GroupList | undefined = await groupsStore.getGroups(request);
      if (!list) {
        error.value = "invalid group list";
        break;
      }
      availableGroups.value = [...availableGroups.value, ...list.groups];
      total = list.total;
      page++;
    }
  } catch (err: unknown) {
    if (err instanceof GatewayError) {
      error.value = err.message;
    } else {
      throw err;
    }
  }
}

async function loadAccountData(username: string | undefined = undefined) {
  const user = username ?? props.account;
  if (!user) {
    return;
  }

  try {
    error.value = undefined;

    const remoteUser: User | undefined = await accountsStore.getAccount(user);

    if (!remoteUser) {
      error.value = "invalid remote user";
      return;
    }

    newUserRequest.value.firstName = remoteUser.firstName;
    newUserRequest.value.lastName = remoteUser.lastName;
    newUserRequest.value.UID = remoteUser.UID;
    newUserRequest.value.GID = remoteUser.GID;
    newUserRequest.value.email = remoteUser.email;
    newUserRequest.value.loginShell = remoteUser.loginShell;
    newUserRequest.value.homeDirectory = remoteUser.homeDirectory;
    newUserRequest.value.username = remoteUser.username;

    const list: GroupList | undefined = await groupsStore.getUserGroups(user);
    if (!list) {
      error.value = "invalid user list";
      return;
    }
    userGroups.value = list.groups;
    userGroupNames.value = list.groups.map((group: Group) => group.name);
  } catch (err: unknown) {
    if (err instanceof GatewayError) {
      error.value = err.message;
    } else {
      throw err;
    }
  }
}

onMounted(async () => {
  error.value = undefined;

  if (!authStore.isAdmin && props.create) {
    error.value = "Log in as admin user to create users";
    return;
  }

  // fetch all available groups
  try {
    checkingGroup.value = true;
    processing.value = true;
    watchGroups.value = false;

    await fetchAvailableGroups();
    if (!props.create) await loadAccountData(props.account);
  } finally {
    processing.value = false;
    checkingGroup.value = false;
    watchGroups.value = true;
  }
});
</script>

<template>
  <div class="account-container">
    <div v-if="error !== undefined">
      <b-alert show variant="danger">
        <h4 class="alert-heading">Error</h4>
        <hr />
        <p class="mb-0">
          <span class="preserve-newlines">
            {{ error }}
          </span>
        </p>
      </b-alert>
    </div>
    <div v-else>
      <b-overlay :show="processing" rounded="sm">
        <b-card
          class="login"
          header-tag="header"
          footer-tag="footer"
          :aria-hidden="processing ? 'true' : null"
        >
          <template v-slot:header>
            <b-row fluid class="text-center">
              <b-col sm="2"></b-col>
              <b-col>{{ title }}</b-col>
              <b-col sm="2"
                ><b-button
                  v-if="!props.create"
                  @click="deleteAccount"
                  pill
                  variant="outline-danger"
                  size="sm"
                  >Delete</b-button
                ></b-col
              >
            </b-row>
          </template>
          <b-card-body>
            <p class="text-left">
              If you are not sure about some values, just leave them blank. You
              will be prompted to enter missing values when you submit.
            </p>
            <b-form autocomplete="off" @submit.prevent="onSubmit">
              <b-form-group
                label-size="sm"
                label-cols-sm="3"
                label="First name:"
                class="account-label"
                label-for="login-input-username"
              >
                <b-form-input
                  autocomplete="off"
                  id="login-input-first-name"
                  size="sm"
                  v-model="newUserRequest.firstName"
                  type="text"
                  required
                  placeholder="Max"
                ></b-form-input>
              </b-form-group>

              <b-form-group
                label-size="sm"
                label-cols-sm="3"
                label="Last name:"
                class="account-label"
                label-for="login-input-last-name"
              >
                <b-form-input
                  autocomplete="off"
                  id="login-input-last-name"
                  size="sm"
                  v-model="newUserRequest.lastName"
                  type="text"
                  required
                  placeholder="Mustermann"
                ></b-form-input>
              </b-form-group>

              <b-form-group
                v-if="all && activeIsAdmin"
                label-size="sm"
                label-cols-sm="3"
                label="UID:"
                class="account-label"
                label-for="login-input-uid"
              >
                <b-form-input
                  autocomplete="off"
                  id="login-input-uid"
                  size="sm"
                  v-model="newUserRequest.UID"
                  type="number"
                  placeholder="2004"
                  aria-describedby="login-input-uid-help-block"
                ></b-form-input>
                <b-form-text id="login-input-uid-help-block">
                  Is optional. If you leave this empty, /bin/bash will be used
                </b-form-text>
              </b-form-group>

              <b-form-group
                v-if="all && activeIsAdmin"
                label-size="sm"
                label-cols-sm="3"
                label="GID:"
                class="account-label"
                label-for="login-input-gid"
              >
                <b-form-input
                  autocomplete="off"
                  id="login-input-gid"
                  size="sm"
                  v-model="newUserRequest.GID"
                  type="number"
                  placeholder="2001"
                  aria-describedby="login-input-gid-help-block"
                ></b-form-input>
                <b-form-text id="login-input-gid-help-block">
                  Is optional. If you leave this empty, /bin/bash will be used
                </b-form-text>
              </b-form-group>

              <b-form-group
                v-if="all"
                label-size="sm"
                label-cols-sm="3"
                label="Login shell:"
                class="account-label"
                label-for="login-input-shell"
              >
                <b-form-input
                  autocomplete="off"
                  id="login-input-shell"
                  size="sm"
                  v-model="newUserRequest.loginShell"
                  type="text"
                  placeholder="/bin/bash"
                  aria-describedby="login-input-shell-help-block"
                ></b-form-input>
                <b-form-text id="login-input-shell-help-block">
                  Is optional. If you leave this empty, /bin/bash will be used
                </b-form-text>
              </b-form-group>

              <b-form-group
                v-if="all"
                label-size="sm"
                label-cols-sm="3"
                label="Home directory:"
                class="account-label"
                label-for="login-input-home-dir"
              >
                <b-form-input
                  autocomplete="off"
                  id="login-input-home-dir"
                  size="sm"
                  v-model="newUserRequest.homeDirectory"
                  type="text"
                  placeholder="/home/max123"
                  aria-describedby="login-input-home-dir-help-block"
                ></b-form-input>
                <b-form-text id="login-input-home-dir-help-block">
                  Is optional. If you leave this empty, the /home/USERNAME is
                  chosen
                </b-form-text>
              </b-form-group>

              <b-form-group
                v-if="activeIsAdmin"
                label-size="sm"
                label-cols-sm="3"
                label="Username:"
                class="account-label"
                label-for="login-input-username"
              >
                <b-form-input
                  autocomplete="off"
                  id="login-input-username"
                  size="sm"
                  v-model="newUserRequest.username"
                  type="text"
                  required
                  placeholder="max123"
                ></b-form-input>
              </b-form-group>

              <b-form-group
                label-size="sm"
                label-cols-sm="3"
                label="Email:"
                class="account-label"
                label-for="login-input-email"
              >
                <b-form-input
                  autocomplete="off"
                  id="login-input-email"
                  size="sm"
                  v-model="newUserRequest.email"
                  :state="validEmail"
                  type="email"
                  required
                  placeholder="max.mustermann@example.com"
                ></b-form-input>
                <b-form-invalid-feedback :state="validEmail">
                  This is not a valid email
                </b-form-invalid-feedback>
                <b-form-valid-feedback :state="validEmail">
                  All good
                </b-form-valid-feedback>
              </b-form-group>

              <b-form-group>
                <b-form-group
                  label-size="sm"
                  label-cols-sm="3"
                  label="Password:"
                  class="account-label"
                  label-for="login-input-password"
                >
                  <b-form-input
                    autocomplete="off"
                    id="login-input-password"
                    size="sm"
                    v-model="newUserRequest.password"
                    type="password"
                    :required="props.create"
                    placeholder=""
                    aria-describedby="login-input-password-help-block"
                  ></b-form-input>
                  <b-form-text id="login-input-password-help-block">
                    Good passwords must be 8-20 characters long, contain letters
                    and numbers, and must not contain spaces, special
                    characters, or emoji.
                  </b-form-text>
                </b-form-group>

                <b-form-group
                  label-size="sm"
                  label-cols-sm="3"
                  label="Confirm:"
                  class="account-label"
                  label-for="login-input-confirm-password"
                >
                  <b-form-input
                    autocomplete="off"
                    id="login-input-confirm-password"
                    size="sm"
                    v-model="passwordConfirm"
                    :state="passwordsMatch"
                    type="password"
                    :required="props.create"
                    placeholder="Confirm password"
                  ></b-form-input>
                  <b-form-invalid-feedback :state="passwordsMatch">
                    Passwords do not match
                  </b-form-invalid-feedback>
                  <b-form-valid-feedback :state="passwordsMatch && newUserRequest.password.length > 0">
                    All good
                  </b-form-valid-feedback>
                </b-form-group>

                <b-row align-h="end">
                  <b-col cols="9"
                    ><b-progress max="6" v-if="enteredPassword">
                      <b-progress-bar
                        :value="passwordStrength"
                        :label="passwordStrengthLabel"
                        :variant="passwordStrengthVariant"
                      ></b-progress-bar></b-progress
                  ></b-col>
                </b-row>
              </b-form-group>

              <b-form-group
                v-if="activeIsAdmin"
                autocomplete="off"
                label-size="sm"
                label-cols-sm="3"
                label="Groups:"
                class="account-label"
                label-for="account-input-groups"
                :state="groupState"
              >
                <b-form-tags
                  v-if="!props.create"
                  autocomplete="off"
                  input-id="account-input-groups"
                  duplicate-tag-text="already in group"
                  :invalid-tag-text="invalidGroupText"
                  :disabled="userGroupInputDisabled"
                  v-model="userGroupNames"
                  tag-variant="primary"
                  tag-pills
                  size="sm"
                  separator=" "
                  :state="groupState"
                  :input-attrs="{ 'aria-describedby': 'tags-validation-help' }"
                  :tag-validator="groupValidator"
                  placeholder="Enter group names separated by spaces"
                  class="mb-2"
                ></b-form-tags>

                <b-form-text v-if="props.create" id="account-input-groups-help">
                  User groups can be edited after the user has been created
                </b-form-text>

                <template v-slot:invalid-feedback>
                  {{ invalidGroupText }}
                </template>
                <template v-slot:description>
                  <div v-if="checkingGroup">Checking group...</div>
                </template>
                <b-alert
                  class="text-left"
                  :show="groupMemberError !== undefined"
                  variant="danger"
                >
                  <span class="preserve-newlines">
                    {{ groupMemberError }}
                  </span>
                </b-alert>
              </b-form-group>

              <b-form-group class="mb-0">
                <b-button
                  class="float-end"
                  size="sm"
                  variant="primary"
                  @click="props.create ? createAccount() : updateAccount()"
                  >{{ create ? "Create account" : "Update" }}
                </b-button>
              </b-form-group>

              <b-alert
                class="text-left"
                :show="submissionError !== undefined"
                variant="danger"
              >
                <span class="preserve-newlines">
                  {{ submissionError }}
                </span>
              </b-alert>
            </b-form>
          </b-card-body>
        </b-card>
      </b-overlay>
    </div>
  </div>
</template>

<style scoped lang="sass">
.account-label
    text-align: right
    font-weight: bold
</style>
