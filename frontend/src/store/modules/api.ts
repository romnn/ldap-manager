import { VuexModule, Module, Action, getModule } from "vuex-module-decorators";
import store from "@/store";
import router from "@/router";
import Vue from "vue";

export const APIEndpoint = "http://localhost:8090/api"

export interface Account {
  uid?: string
  givenname?: string
  sn?: string
  mail?: string
  [key: string]: string | undefined
}

export interface APIState {

}

@Module({ dynamic: true, store, name: "api" })
class APIMod extends VuexModule implements APIState {

  get filename(): (name: string) => string {
    return (name: string) => name.replace(/[^a-z0-9]/gi, "_").toLowerCase();
  }

  @Action({ rawError: true })
  public async listAccounts(): Promise<Account[]> {
    return new Promise<Account[]>((resolve, reject) => {
      Vue.axios.post(APIEndpoint + "/accounts/list", {}).then(
        response => {
          debugger;
          // Get token
          let token = response.data.token;
          localStorage.setItem("user-token", token);
          Vue.axios.defaults.headers.common["Authorization"] = token;
          commit("AUTH_SUCCESS", token);
          // dispatch("USER_REQUEST");
          resolve(response);
        },
        err => {
          debugger;
          // commit("AUTH_ERROR", err);
          localStorage.removeItem("user-token");
          if (err.response) {
            let error_response = err.response.data;
            reject({
              status: err.response.status,
              error: error_response.error,
              message: error_response.message
            });
          }
          // Reject promise
          reject({
            status: null,
            error: "Authentication failed",
            message:
              "Could not properly connect to the server. Please try again later."
          });
        }
      );
      /*
      const descr = new StockItemsDescriptor();
      descr.setType(Stock.Type.BEAT);
      descr.setFilter(filter);

      const sl = new StockList();
      const stream = this.client.getStockItems(descr, undefined);
      stream.on("data", (item: Stock) => {
        sl.addItems(item);
      });
      stream.on("error", (err: Error) => {
        reject(err);
      });
      stream.on("status", (status: Status) => {
        const total = status.metadata?.total;
        if (total != undefined) sl.setTotal(parseInt(total));
        const limit = status.metadata?.limit;
        if (limit != undefined) sl.setLimit(parseInt(limit));
      });
      stream.on("end", () => {
        resolve(sl);
      });
      */
    });
  }

  /*
  get stockIndex(): (stockID?: Stock.ID) => string | undefined {
    // Stringified for use with indices
    return (stockID?: Stock.ID) => {
      if (stockID == undefined) return undefined;
      const type = Object.keys(Stock.Type)[stockID?.getType()]?.toUpperCase();
      const id = stockID?.getId().toString();
      if (type !== undefined && id !== undefined) return type + id;
      return undefined;
    };
  }

  get licenseName(): (type: License.Type) => string {
    return (type: License.Type) => Object.keys(License.Type)[type];
  }

  get previewAudio(): (s: Stock.ID) => string {
    return (s: Stock.ID) => "/static/previews/" + this.stockIndex(s) + ".mp3";
  }

  get stockCover(): (s: Stock.ID) => string {
    return (s: Stock.ID) => "/static/artwork/" + this.stockIndex(s);
  }
  get redirect(): (platform: string) => string {
    return (platform: string): string => {
      const loc: Location = {
        path: "/goto/" + platform,
        params: {
          return: router.currentRoute.fullPath
        }
      };
      return router.resolve(loc).href;
    };
  }
  get redirectBeatTo(): (
    source: string,
    id: string,
    platform: string
  ) => string {
    return (source: string, id: string, platform: string): string => {
      const loc: Location = {
        path: "/goto/" + source + "/" + id + "/" + platform,
        params: {
          return: router.currentRoute.fullPath
        }
      };
      return router.resolve(loc).href;
    };
  }

  @Action({ rawError: true })
  public async loadStock(id: Stock.ID): Promise<Stock> {
    const descr = new StockItemDescriptor();
    descr.setId(id);
    return this.client.getStockItem(descr, null);
  }

  @Action({ rawError: true })
  public async loadRelated(id: Stock.ID): Promise<StockList> {
    return new Promise<StockList>((resolve, reject) => {
      const descr = new StockItemDescriptor();
      descr.setId(id);
      const sl = new StockList();
      const stream = this.client.getRelatedStockItems(descr, undefined);
      stream.on("data", (item: Stock) => {
        sl.addItems(item);
      });
      stream.on("status", () => {
        // Noop
      });
      stream.on("error", (err: Error) => {
        reject(err);
      });
      stream.on("end", () => {
        resolve(sl);
      });
    });
  }
  */
}

export const APIModule = getModule(APIMod);
