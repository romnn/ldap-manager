<script setup lang="ts">
const emit = defineEmits(["confirm", "cancel"]);

const props = withDefaults(
  defineProps<{
    message: string;
    ackMessage: string;
    processing?: boolean;
  }>(),
  {
    message: "Are you sure?",
    ackMessage: "Yes",
  }
);

function confirm() {
  emit("confirm");
}

function cancel() {
  emit("cancel");
}
</script>

<template>
  <div class="confirmation">
    <b-overlay :show="props.processing" rounded="sm">
      <span class="m-2">{{ props.message }}</span>
      <br />
      <span>
        <b-button size="sm" class="m-2" @click="cancel()" variant="light"
          >Cancel</b-button
        >
        <b-button size="sm" class="m-2" @click="confirm()" variant="danger">{{
          props.ackMessage
        }}</b-button>
      </span>
    </b-overlay>
  </div>
</template>

<style lang="sass">
.confirmation
  border: 2px #e9ecef solid
  border-radius: 15px
  padding: 15px
  background-color: #ffffff
  z-index: 999999
  position: fixed
  top: 50%
  left: 50%
  transform: translate(-50%, -50%)
</style>
