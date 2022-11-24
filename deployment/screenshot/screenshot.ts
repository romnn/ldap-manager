import { status } from "@grpc/grpc-js";
import axios, { AxiosError } from "axios";
import {
  blue,
  Color,
  cyan,
  green,
  magenta,
  red,
  gray,
  white,
  yellow,
} from "colorette";
import { LoginRequest, NewUserRequest, Token } from "ldap-manager";
import path from "path";
import puppeteer from "puppeteer";
import { PageEmittedEvents, PageEventObject } from "puppeteer";

const ADMIN_USERNAME = process.env.ADMIN_USERNAME ?? "ldapadmin";
const ADMIN_PASSWORD = process.env.ADMIN_PASSWORD ?? "changeme";
const ADDRESS = `http://${process.env.LDAP_MANAGER_HOST ?? "localhost"}`;
const OUTPUT_DIR =
  process.env.OUTPUT_DIR ?? path.join(__dirname, "../../screenshots");

(async () => {
  // get admin authentication token
  let adminToken = await authenticate({ username: "admin", password: "admin" });

  // create a new user
  const request: NewUserRequest = {
    username: "test-user",
    email: "somebody@mail.com",
    firstName: "Some",
    lastName: "User",
    password: "changeme",
    loginShell: "",
    homeDirectory: "",
    UID: 0,
    GID: 0,
  };
  await newAccount(request, { token: adminToken });
  let userToken = await authenticate({
    username: "test-user",
    password: "changeme",
  });

  await screenshotAll({ adminToken, userToken });
})();

async function newAccount(
  request: NewUserRequest,
  { token }: { token: Token | undefined }
) {
  const response = await axios.put(ADDRESS + "/api/v1/user", request, {
    headers: {
      "x-user-token": token?.token ?? "",
    },
    validateStatus: (status: number) => true,
  });
  if (response.status !== 200 && response.data.code !== status.ALREADY_EXISTS) {
    throw response.data;
  }
}

async function authenticate(request: LoginRequest): Promise<Token | undefined> {
  const response = await axios.post(ADDRESS + "/api/v1/login", request, {
    validateStatus: (status: number) => true,
  });
  if (response.status != 200) {
    throw response.data;
  }
  return Token.fromJSON(response.data);
}

async function screenshotAll({
  adminToken,
  userToken,
}: {
  adminToken: Token | undefined;
  userToken: Token | undefined;
}) {
  const browser = await puppeteer.launch({
    // bindAddress : '0.0.0.0',
    args: [
      "--headless",
      "--disable-gpu",
      "--disable-dev-shm-usage",
      "--remote-debugging-port=9222",
      "--remote-debugging-address=0.0.0.0",
    ],
  });

  let pages: { token?: Token; url: string; filename: string }[] = [
    { url: ADDRESS + "/login", filename: "login" },
  ];

  let users: { [name: string]: Token | undefined } = {
    admin: adminToken,
    user: userToken,
  };

  for (const user in users) {
    const token = users[user];
    const username = token?.username;
    pages.push({ token, url: ADDRESS, filename: `home-of-${user}` });
    pages.push({
      token,
      url: `${ADDRESS}/users/list`,
      filename: `users-list-for-${user}`,
    });
    if (username) {
      pages.push({
        token,
        url: `${ADDRESS}/user/${username}`,
        filename: `user-edit-by-${user}`,
      });
    }
    pages.push({
      token,
      url: `${ADDRESS}/groups/list`,
      filename: `groups-list-for-${user}`,
    });
    pages.push({
      token,
      url: `${ADDRESS}/group/users`,
      filename: `group-edit-by-${user}`,
    });
  }
  await Promise.all(
    pages.map(async (page) => {
      let message = `taking screenshot of ${page.url}`;
      if (page.token !== undefined)
        message += ` for user ${page.token.username}`;
      console.log(cyan(message));
      await screenshotPage({
        browser,
        token: page.token,
        url: page.url,
        filename: page.filename,
      });
    })
  );
  console.log(green("done"));
  await browser.close();
}

async function screenshotPage({
  browser,
  token,
  url,
  filename,
}: {
  browser: any;
  token?: Token;
  url: string;
  filename: string;
}) {
  const context = await browser.createIncognitoBrowserContext();
  const page = await context.newPage();
  await page.setCacheEnabled(false);
  await page.setViewport({
    width: 900,
    height: 900,
    deviceScaleFactor: 3,
  });

  // https://github.com/puppeteer/puppeteer/blob/main/packages/puppeteer-core/src/api/Page.ts#L462
  page.on(
    PageEmittedEvents.Console,
    (event: PageEventObject[PageEmittedEvents.Console]) => {
      const typ = event.type().substr(0, 3).toUpperCase();
      const colors: { [key: string]: Color } = {
        ERR: red,
        WAR: yellow,
        INF: cyan,
      };
      const color = colors[typ] || white;
      console.log(color(`${typ} ${event.text()}`));
    }
  );
  page.on(
    PageEmittedEvents.PageError,
    ({ message }: PageEventObject[PageEmittedEvents.PageError]) =>
      console.log(red(message))
  );
  if (process.env.VERBOSE) {
    page.on(
      PageEmittedEvents.Response,
      (response: PageEventObject[PageEmittedEvents.Response]) => {
        console.log(gray(`${response.status()} ${response.url()}`));
      }
    );
  }
  page.on(
    PageEmittedEvents.RequestFailed,
    (request: PageEventObject[PageEmittedEvents.RequestFailed]) =>
      console.log(red(`${request.failure()?.errorText} ${request.url()}`))
  );

  if (token) {
    await page.goto(url, { waitUntil: "networkidle2" });
    await page.evaluate((token: Token) => {
      const AUTH_TOKEN_KEY = "auth/token";
      const AUTH_IS_ADMIN_KEY = "auth/admin";
      const AUTH_USERNAME_KEY = "auth/username";
      const AUTH_DISPLAY_NAME_KEY = "auth/displayname";

      // @ts-ignore
      localStorage.setItem(AUTH_TOKEN_KEY, token.token);
      // @ts-ignore
      localStorage.setItem(AUTH_USERNAME_KEY, token.username);
      if (token.isAdmin) {
        // @ts-ignore
        localStorage.setItem(AUTH_IS_ADMIN_KEY, "true");
      }
      // @ts-ignore
      localStorage.setItem(AUTH_DISPLAY_NAME_KEY, token.displayName);
    }, token);
  }
  await page.goto(url, { waitUntil: "networkidle2" });
  await page.screenshot({ path: path.join(OUTPUT_DIR, `${filename}.png`) });
}
