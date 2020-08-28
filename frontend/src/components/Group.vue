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
          <b-form @submit.prevent="onSubmit">
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
                v-model="form.gid"
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
                <b-col><member-list-c></member-list-c></b-col>
                <b-col><member-list-c></member-list-c></b-col>
              </b-row>
            </b-form-group>

            <b-form-group>
              <b-button
                class="float-right"
                size="sm"
                type="submit"
                variant="primary"
                >{{ create ? "Create account" : "Update" }}
              </b-button>
            </b-form-group>
          </b-form>
        </b-card-body>
      </b-card>
    </b-overlay>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
import { GatewayError } from "../types";
import { GroupModule } from "../store/modules/groups";
import { AppModule } from "../store/modules/app";
import MemberListC from "./MemberList.vue";

@Component({
  components: {MemberListC}
})
export default class GroupC extends Vue {
  @Prop() private name!: string;
  @Prop({ default: "Group" }) private title!: string;
  @Prop({ default: false }) private all!: boolean;
  @Prop({ default: false }) private create!: boolean;

  protected processing = false;

  protected form = {
    name: "",
    gid: 0
  };

  deleteGroup(name: string) {
    AppModule.newConfirmation({ message: "Are you sure?", ack: "Yes, delete" })
      .then(() => {
        this.processing = true;
        GroupModule.deleteGroup(name)
          .then(() => this.$router.push({ name: "GroupsRoute" }))
          .catch((err: GatewayError) => alert(err.message))
          .finally(() => (this.processing = false));
      })
      .catch(() => {
        // Ingore
      });
  }

  createGroup() {
    this.processing = true;
    GroupModule.newGroup(this.form)
      .catch((err: GatewayError) => alert(err.message))
      .finally(() => (this.processing = false));
  }

  updateGroup() {
    this.processing = true;
  }

  onSubmit() {
    this.create ? this.createGroup() : this.updateGroup();
  }
}
</script>

<style scoped lang="sass">
group-label
    text-align: right
    font-weight: bold
</style>
