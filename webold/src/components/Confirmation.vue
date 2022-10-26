<template>
  <div class="confirmation">
    <b-overlay :show="processing" rounded="sm">
      <span class="m-2">{{ message }}</span>
      <br />
      <span>
        <b-button size="sm" class="m-2" @click="cancel()" variant="light"
          >Cancel</b-button
        >
        <b-button size="sm" class="m-2" @click="confirm()" variant="danger">{{
          ackMessage
        }}</b-button>
      </span>
    </b-overlay>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";

@Component
export default class ConfirmationC extends Vue {
  @Prop({ default: "Are you sure?" }) private message!: string;
  @Prop({ default: "Yes" }) private ackMessage!: string;
  @Prop() private processing!: boolean;

  cancel() {
    this.$emit("cancel");
  }

  confirm() {
    this.$emit("confirm");
  }
}
</script>

<style lang="sass">
.confirmation
  border: 2px #e9ecef solid
  border-radius: 15px
  padding: 15px
  background-color: #ffffff
  z-index: 999999
  position: fixed
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
</style>
