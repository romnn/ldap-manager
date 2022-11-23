<script setup lang="ts">
import { ref, defineProps, computed, onMounted } from "vue";
import { GatewayError } from "../constants";
import type { Group, UserList } from "ldap-manager";

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
const accountStore = useAccountsStore();
const groupStore = useGroupsStore();
const memberStore = useMembersStore();

const search = ref("");
const processing = ref(false);
const loadingMembers = ref(false);
const loadingAvailableAccounts = ref(false);

const loadingAvailableError = ref<string | undefined>(undefined);
const loadingGroupError = ref<string | undefined>(undefined);
const groupMemberOperationError = ref<string | undefined>(undefined);
const submissionError = ref<string | undefined>(undefined);

const available = ref<UserList>({ users: [], total: 0 });

const membersSearch = ref<string>("");
const availableSearch = ref<string>("");

const form = ref< {
  members: string[];
  name: string;
  GID: number;
}>({
  members: [],
  name: "",
  GID: 0,
});

const props = withDefaults(defineProps<{
   name?: string
   title: string,
   all: boolean,
   create: boolean,
 }>(), {
     title: 'Group',
     all: false,
     create: false,
 });

const filteredMembers = computed(() =>
  form.value.members.filter((member) => {
    return member.includes(membersSearch.value);
  })
);

async function updateAvailableSearch(search: string) {
  availableSearch.value = search;
  await loadAvailableAccounts();
}

function isMember(username: string) {
  return form.value.members.includes(username);
}

function successAlert(message: string) {
  toast?.success({
    title: "Success",
    body: message,
  }, {
    autoHide: true,
    delay: 5000,
  });
}

async function deleteGroup(name: string | undefined) {
  const groupName = name ?? props.name;
  if (!groupName) {
    return;
  }
  try {
    await appStore.newConfirmation({ message: "Are you sure?", ack: "Yes, delete" });
  } catch (err: unknown) {
    return;
  }
  try {
    processing.value = true;
    await groupStore.deleteGroup(groupName);
    /* this.$router.push({ name: "GroupsRoute" }); */
    successAlert(`${groupName} was deleted`);
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

async function createGroup() {
  processing.value = true;
  try {
    const newGroupRequest = form.value;
    console.log(newGroupRequest );
    await groupStore.newGroup(newGroupRequest);
    successAlert(`${newGroupRequest.name} was created`);
    router
      .push({
        name: "EditGroupRoute",
        params: { name: newGroupRequest.name }
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

async function removeAccount(username: string, group: string | undefined = undefined) {
  if (props.create) {
    form.value.members = form.value.members.filter(
      (member) => member !== username
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
      username: username,
      group: groupName,
    });
    successAlert(`${username} was removed from ${groupName}`);
    form.value.members = form.value.members.filter(
      member => member !== username
    );
  } catch (err: unknown) {
    if (err instanceof GatewayError) {
      groupMemberOperationError.value = err.message;
    } else {
      throw err;
    }
  } finally {
    processing.value = false;
  }
}

async function addAccount(username: string, group: string | undefined = undefined) {
  const groupName = group ?? props.name;
  if (!groupName) {
    return;
  }
  if (props.create) {
    form.value.members.push(username);
    return;
  }
  processing.value = true;
  groupMemberOperationError.value = undefined;
  try {
    await memberStore.addGroupMember({
      username: username,
      group: groupName,
    });
    successAlert(`${username} was added to ${groupName}`);
    form.value.members.push(username);
  } catch (err: unknown) {
    if (err instanceof GatewayError) {
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
    router
      .push({
        name: "EditGroupRoute",
        params: { name: form.value.name }
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

async function loadAvailableAccounts() {
  try {
    loadingAvailableAccounts.value = true;
    loadingAvailableError.value = undefined;
    available.value = { users: [], total: 0 };

    const list: UserList | undefined = await accountStore.listAccounts({
      search: availableSearch.value,
      page: 1,
      perPage: 50,
    });
    if (!list) {
      loadingAvailableError.value = "invalid user list";
      return;
    }
    available.value.users = list?.users ?? [];
    available.value.total = list?.total ?? "0";
  } catch (err: unknown) {
    if (err instanceof GatewayError) {
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
    form.value.members = group.members;
  } catch (err: unknown) {
    if (err instanceof GatewayError) {
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
                v-if="!create"
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
              v-if="all"
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
                <b-col>
                  <!--
                  <member-list-c
                    title="Members"
                    :loading="loadingMembers"
                    v-on:search="updateMemberSearch"
                  >
                    <div v-if="form.members.length < 1">
                      No members yet
                    </div>
                    <table v-else class="striped-table">
                      <thead>
                        <td>Username</td>
                        <td></td>
                      </thead>
                      <tr
                        v-for="(member, idx) in filteredMembers"
                        v-bind:key="member"
                        :class="{
                          even: idx % 2 == 0
                        }"
                      >
                        <td>{{ member }}</td>
                        <td>
                          <b-button
                            pill
                            @click="removeAccount(member)"
                            size="sm"
                            class="mr-2 float-right"
                            variant="outline-danger"
                            >Remove</b-button
                          >
                        </td>
                      </tr>
                    </table>
                  </member-list-c>
                  -->
                </b-col>
                <b-col>
                  <!--
                  <member-list-c
                    title="All users"
                    :loading="loadingAvailableAccounts"
                    v-on:search="updateAvailableSearch"
                  >
                    <div v-if="available.users.length < 1">
                      No users available
                    </div>
                    <table v-else class="striped-table">
                      <thead>
                        <td>Username</td>
                        <td></td>
                      </thead>
                      <tr
                        v-for="(user, idx) in available.users"
                        v-bind:key="user.data.uid"
                        :class="{
                          even: idx % 2 == 0,
                          isMember: isMember(user.data.uid)
                        }"
                      >
                        <td>{{ user.data.uid }}</td>
                        <td>
                          <span v-if="isMember(user.data.uid)">
                            <i>member already</i>
                          </span>
                          <div v-else>
                            <b-button
                              pill
                              @click="addAccount(user.data.uid)"
                              size="sm"
                              class="mr-2 float-right"
                              variant="outline-primary"
                              >Add</b-button
                            >
                          </div>
                        </td>
                      </tr>
                    </table>
                  </member-list-c>
                  -->
                  <b-alert
                    class="text-left"
                    dismissible
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
              dismissible
              :show="loadingGroupError !== undefined"
              variant="danger"
            >
              {{ loadingGroupError }}
            </b-alert>

            <b-alert
              class="text-left"
              dismissible
              :show="groupMemberOperationError !== undefined"
              variant="danger"
            >
              {{ groupMemberOperationError }}
            </b-alert>

            <b-form-group>
              <b-button
                class="float-right"
                size="sm"
                variant="primary"
                @click="props.create ? createGroup() : updateGroup()"
                >{{ props.create ? "Create group" : "Update" }}
              </b-button>
            </b-form-group>

            <b-alert
              class="text-left"
              dismissible
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
