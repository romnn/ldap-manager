import {
  VuexModule,
  Module,
  Mutation,
  Action,
  getModule
} from "vuex-module-decorators";
import store from "@/store";
import router from "@/router";

export interface PendingConfirmation {
  message: string;
  ack: string;
  resolve?: () => void;
  reject?: () => void;
}

export interface AppState {
  pendingConfirmation: PendingConfirmation | null;
}

@Module({ dynamic: true, store, name: "app" })
class AppMod extends VuexModule implements AppState {
  allowAll = true;
  pendingConfirmation: PendingConfirmation | null = null;

  @Mutation
  public setPendingConfirmation(p: PendingConfirmation) {
    this.pendingConfirmation = p;
  }

  @Action({ rawError: true })
  public newConfirmation(req: { message: string; ack: string }): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      this.setPendingConfirmation({
        message: req.message,
        ack: req.ack,
        resolve: resolve,
        reject: reject
      });
    });
  }

  @Mutation
  public cancelConfirmation() {
    if (this.pendingConfirmation && this.pendingConfirmation.reject)
      this.pendingConfirmation.reject();
    this.pendingConfirmation = null;
  }

  @Mutation
  public confirmConfirmation() {
    if (this.pendingConfirmation && this.pendingConfirmation.reject)
      this.pendingConfirmation.reject();
    this.pendingConfirmation = null;
  }

  get all() {
    return (
      this.allowAll &&
      Object.keys(router.currentRoute?.query ?? {}).includes("all")
    );
  }
}

export const AppModule = getModule(AppMod);
