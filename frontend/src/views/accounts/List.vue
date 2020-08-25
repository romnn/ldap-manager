<template>
  <div class="account-list-container">
    <b-card body-class="p-1" bg-variant="dark" v-for="user in list">
      <div>
        This is some content without the default
        <samp>&lt;b-card-body&gt;</samp> section. Notice the lack of padding
        between the card's border and this gray <samp>&lt;div&gt;</samp>.
      </div>
    </b-card>
    <!--
    <b-card class="login" header-tag="header" footer-tag="footer">
      <template v-slot:header>
        <h6 class="mb-0">Login</h6>
      </template>
      <b-card-body>
        <b-form @submit="onSubmit" @reset="onReset">
          <b-form-group
            label-size="sm"
            label-cols-sm="4"
            label-cols-lg="3"
            label="Username:"
            label-for="login-input-username"
          >
            <b-form-input
              id="login-input-username"
              size="sm"
              v-model="form.username"
              type="text"
              required
              placeholder="Enter username"
            ></b-form-input>
          </b-form-group>

          <b-form-group
            label-size="sm"
            label-cols-sm="4"
            label-cols-lg="3"
            label="Password:"
            label-for="login-input-password"
          >
            <b-form-input
              id="login-input-password"
              size="sm"
              v-model="form.password"
              type="password"
              required
              placeholder="Enter password"
            ></b-form-input>
          </b-form-group>

          <b-form-group
            inline
            label-cols-sm="3"
            class="mb-0"
          >
            <b-form-row>
            <b-col>
              <b-form-checkbox size="sm" v-model="form.remember">Remember me</b-form-checkbox>
            </b-col>
            <b-col>
              <b-button size="sm" type="submit" variant="primary">Log in</b-button>
            </b-col>
          </b-form-row>
          </b-form-group>
        </b-form>
      </b-card-body>
    </b-card>
    -->
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { AccountModule, UserList } from "../../store/modules/accounts";

@Component({
  components: {}
})
export default class AccountList extends Vue {
  list: UserList = { users: [] };

  mounted() {
    AccountModule.listAccounts()
      .then((list: UserList) => {
        this.list = list;
      })
      .catch(err => {
        debugger;
      });
  }
}
</script>

<style lang="sass" scoped></style>
