<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { GatewayError } from "../constants";
import { AxiosError } from "axios";
import MemberListComponent from "./MemberListComponent.vue";
import type { Group, User, UserList } from "ldap-manager";

import { useRouter } from "vue-router";
import { useToast } from "bootstrap-vue-3";
import { useGroupsStore } from "../stores/groups";
import { useMembersStore } from "../stores/members";
import { useAppStore } from "../stores/app";
import { useAccountsStore } from "../stores/accounts";

const router = useRouter();
const toast = useToast();
const appStore = useAppStore();
const accountStore = useAccountsStore();
const groupStore = useGroupsStore();
const memberStore = useMembersStore();

const processing = ref(false);
const loadingMembers = ref(false);
const loadingAvailableAccounts = ref(false);

const loadingAvailableError = ref<string | undefined>(undefined);
const loadingGroupError = ref<string | undefined>(undefined);
const groupMemberOperationError = ref<string | undefined>(undefined);
const submissionError = ref<string | undefined>(undefined);

const availableMap = ref<Map<string, User>>(new Map());

const memberSearch = ref<string>("");
const availableSearch = ref<string>("");

const form = ref<{
  /* members: string[]; */
  members: User[];
  name: string;
  GID: number;
}>({
  members: [],
  name: "",
  GID: 0,
});

const props = withDefaults(
  defineProps<{
    name?: string;
    title?: string;
    internal?: boolean;
    create?: boolean;
  }>(),
  {
    title: "Group",
    internal: false,
    create: false,
  }
);

const members = computed((): User[] =>
  form.value.members
    .map((member) => {
      // matching DN from available users
      return availableMap.value.get(member.DN);
    })
    .filter((member): member is User => !!member)
);

const filteredMembers = computed((): User[] =>
  members.value.filter((member) => {
    return member.username.includes(memberSearch.value);
  })
);

async function updateAvailableSearch(search: string) {
  availableSearch.value = search;
  await loadAvailableAccounts();
}

function updateMemberSearch(search: string) {
  memberSearch.value = search;
}

function isMember(user: User) {
  return form.value.members.some((member) => member.DN === user.DN);
}

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

async function deleteGroup(name: string | undefined) {
  const groupName = name ?? props.name;
  if (!groupName) {
    return;
  }
  try {
    await appStore.newConfirmation({
      message: "Are you sure?",
      ack: "Yes, delete",
    });
  } catch (err: unknown) {
    return;
  }
  try {
    processing.value = true;
    await groupStore.deleteGroup(groupName);
    successAlert(`${groupName} was deleted`);
    router.push({ name: "GroupsRoute" });
  } catch (err: unknown) {
    if (err instanceof GatewayError) {
      submissionError.value = err.message;
    } else if (err instanceof AxiosError) {
      submissionError.value = err.message;
    } else {
      throw err;
    }
  } finally {
    processing.value = false;
  }
}

async function createGroup() {
  processing.value = true;
  try {
    const newGroupRequest = {
      ...form.value,
      members: form.value.members.map((member) => member.username),
    };
    await groupStore.newGroup(newGroupRequest);
    successAlert(`${newGroupRequest.name} was created`);
    router.push({
      name: "EditGroupRoute",
      params: { name: newGroupRequest.name },
    });
  } catch (err: unknown) {
    if (err instanceof GatewayError) {
      submissionError.value = err.message;
    } else if (err instanceof AxiosError) {
      submissionError.value = err.message;
    } else {
      throw err;
    }
  } finally {
    processing.value = false;
  }
}

async function removeAccount(
  user: User,
  group: string | undefined = undefined
) {
  if (props.create) {
    form.value.members = form.value.members.filter(
      (member) => member.DN !== user.DN
    );
    return;
  }
  const groupName = group ?? props.name;
  if (!groupName) {
    return;
  }

  processing.value = true;
  groupMemberOperationError.value = undefined;

  try {
    await memberStore.removeGroupMember({
      username: user.username,
      group: groupName,
      dn: "", // ignored
    });
    successAlert(`${user.username} was removed from ${groupName}`);
    form.value.members = form.value.members.filter(
      (member) => member.DN !== user.DN
    );
  } catch (err: unknown) {
    if (err instanceof GatewayError) {
      groupMemberOperationError.value = err.message;
    } else if (err instanceof AxiosError) {
      groupMemberOperationError.value = err.message;
    } else {
      throw err;
    }
  } finally {
    processing.value = false;
  }
}

async function addAccount(user: User, group: string | undefined = undefined) {
  if (props.create) {
    form.value.members.push(user);
    console.log(form.value.members);
    return;
  }

  const groupName = group ?? props.name;
  if (!groupName) {
    return;
  }
  processing.value = true;
  groupMemberOperationError.value = undefined;
  try {
    await memberStore.addGroupMember({
      username: user.username,
      group: groupName,
      dn: "", // ignored
    });
    successAlert(`${user.username} was added to ${groupName}`);
    form.value.members.push(user);
  } catch (err: unknown) {
    if (err instanceof GatewayError) {
      groupMemberOperationError.value = err.message;
    } else if (err instanceof AxiosError) {
      groupMemberOperationError.value = err.message;
    } else {
      throw err;
    }
  } finally {
    processing.value = false;
  }
}

async function updateGroup(group: string | undefined = undefined) {
  const oldGroupName = group ?? props.name;
  if (!oldGroupName) {
    return;
  }

  processing.value = true;
  try {
    await groupStore.updateGroup({
      name: oldGroupName,
      newName: form.value.name,
      GID: form.value.GID,
    });
    successAlert(`${oldGroupName} was updated`);
    router.push({
      name: "EditGroupRoute",
      params: { name: form.value.name },
    });
  } catch (err: unknown) {
    if (err instanceof GatewayError) {
      submissionError.value = err.message;
    } else if (err instanceof AxiosError) {
      submissionError.value = err.message;
    } else {
      throw err;
    }
  } finally {
    processing.value = false;
  }
}

async function loadAvailableAccounts() {
  try {
    loadingAvailableAccounts.value = true;
    loadingAvailableError.value = undefined;
    availableMap.value = new Map();

    const list: UserList | undefined = await accountStore.listAccounts({
      search: availableSearch.value,
      page: 1,
      perPage: 50,
    });
    if (!list) {
      loadingAvailableError.value = "invalid user list";
      return;
    }
    console.log(list.users);
    for (const user of list.users) {
      availableMap.value.set(user.DN, user);
    }
  } catch (err: unknown) {
    if (err instanceof GatewayError) {
      loadingAvailableError.value = err.message;
    } else if (err instanceof AxiosError) {
      loadingAvailableError.value = err.message;
    } else {
      throw err;
    }
  } finally {
    loadingAvailableAccounts.value = false;
  }
}

async function loadGroupData(name: string | undefined = undefined) {
  try {
    const groupName = name ?? props.name;
    if (!groupName) {
      return;
    }
    loadingMembers.value = true;
    loadingGroupError.value = undefined;

    const group: Group | undefined = await groupStore.getGroup(groupName);
    if (!group) {
      loadingGroupError.value = "invalid group data";
      return;
    }
    form.value.GID = group.GID;
    form.value.name = group.name;
    // todo: find them in the available users
    form.value.members = group.members
      .map((member) => availableMap.value.get(member.dn))
      .filter((member): member is User => !!member);
    console.log(group.members);
  } catch (err: unknown) {
    if (err instanceof GatewayError) {
      loadingGroupError.value = err.message;
    } else if (err instanceof AxiosError) {
      loadingGroupError.value = err.message;
    } else {
      throw err;
    }
  } finally {
    loadingMembers.value = false;
  }
}

onMounted(async () => {
  await loadAvailableAccounts();
  // load group data after available accounts
  if (!props.create) await loadGroupData();
});
</script>

<template>
  <div class="account-container">
    <b-overlay :show="processing" rounded="sm">
      <b-card
        class="login"
        header-tag="header"
        footer-tag="footer"
        :aria-hidden="processing ? 'true' : null"
      >
        <template v-slot:header>
          <b-row class="text-center">
            <b-col></b-col>
            <b-col cols="8">{{ title }}</b-col>
            <b-col
              ><b-button
                v-if="!props.create"
                @click="deleteGroup"
                pill
                variant="outline-danger"
                size="sm"
                >Delete</b-button
              ></b-col
            >
          </b-row>
        </template>
        <b-card-body>
          <b-form>
            <b-form-group
              label-size="sm"
              label-cols-sm="3"
              label="Group name:"
              class="group-label"
              label-for="login-input-name"
            >
              <b-form-input
                autocomplete="off"
                id="login-input-name"
                size="sm"
                v-model="form.name"
                type="text"
                required
                placeholder="My group"
              ></b-form-input>
            </b-form-group>

            <b-form-group
              v-if="props.internal"
              label-size="sm"
              label-cols-sm="3"
              label="GID:"
              class="group-label"
              label-for="group-input-gid"
            >
              <b-form-input
                autocomplete="off"
                id="group-input-gid"
                size="sm"
                v-model="form.GID"
                type="number"
                placeholder="2001"
                aria-describedby="group-input-gid-help-block"
              ></b-form-input>
              <b-form-text id="group-input-gid-help-block">
                Is optional. If you leave this empty, will be auto calculated
              </b-form-text>
            </b-form-group>

            <b-form-group>
              <b-row>
                <b-col :style="{ overflow: 'hidden' }">
                  <member-list-component
                    class="member-list"
                    title="Members"
                    :loading="loadingMembers"
                    v-on:search="updateMemberSearch"
                  >
                    <div v-if="form.members.length < 1">No members yet</div>
                    <table v-else class="striped-table">
                      <thead>
                        <td>Username</td>
                        <td></td>
                      </thead>
                      <tr
                        v-for="(member, idx) in filteredMembers"
                        v-bind:key="member.UID"
                        :class="{
                          even: idx % 2 == 0,
                        }"
                      >
                        <td>{{ member.username }}</td>
                        <td>
                          <b-button
                            pill
                            @click="removeAccount(member)"
                            size="sm"
                            class="mr-2 float-end"
                            variant="outline-danger"
                            >Remove</b-button
                          >
                        </td>
                      </tr>
                    </table>
                  </member-list-component>
                </b-col>
                <b-col>
                  <member-list-component
                    title="All users"
                    :loading="loadingAvailableAccounts"
                    v-on:search="updateAvailableSearch"
                  >
                    <div v-if="availableMap.size < 1">No users available</div>
                    <table v-else class="striped-table">
                      <thead>
                        <td>Username</td>
                        <td></td>
                      </thead>
                      <tr
                        v-for="(user, idx) in availableMap.values()"
                        v-bind:key="user.UID"
                        :class="{
                          even: idx % 2 == 0,
                          isMember: isMember(user),
                        }"
                      >
                        <td>{{ user.username }}</td>
                        <td>
                          <span v-if="isMember(user)">
                            <i>member already</i>
                          </span>
                          <div v-else>
                            <b-button
                              pill
                              @click="addAccount(user)"
                              size="sm"
                              class="mr-2 float-end"
                              variant="outline-primary"
                              >Add</b-button
                            >
                          </div>
                        </td>
                      </tr>
                    </table>
                  </member-list-component>
                  <b-alert
                    class="text-left"
                    :show="loadingAvailableError !== undefined"
                    variant="danger"
                  >
                    {{ loadingAvailableError }}
                  </b-alert>
                </b-col>
              </b-row>
            </b-form-group>

            <b-alert
              class="text-left"
              :show="loadingGroupError !== undefined"
              variant="danger"
            >
              {{ loadingGroupError }}
            </b-alert>

            <b-alert
              class="text-left"
              :show="groupMemberOperationError !== undefined"
              variant="danger"
            >
              {{ groupMemberOperationError }}
            </b-alert>

            <b-form-group>
              <b-button
                class="float-end"
                size="sm"
                variant="primary"
                @click="props.create ? createGroup() : updateGroup()"
                >{{ props.create ? "Create group" : "Update" }}
              </b-button>
            </b-form-group>

            <b-alert
              class="text-left"
              :show="submissionError !== undefined"
              variant="danger"
            >
              {{ submissionError }}
            </b-alert>
          </b-form>
        </b-card-body>
      </b-card>
    </b-overlay>
  </div>
</template>

<style scoped lang="sass">
group-label
    text-align: right
    font-weight: bold
</style>
