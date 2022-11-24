import axios, { AxiosError } from 'axios';
import {Token} from 'ldap-manager';
import puppeteer from 'puppeteer';

const ADDRESS = 'http://ldapmanager';

const ALREADY_EXISTS = 6;

(async () => {
  // get admin authentication token
  let admin = await authenticate({username: 'ldapadmin', password: 'changeme'});

  // create a new user
  // const request = {};
  // await newAccount(request);
  // token = admin.token, username = 'my-user', password = 'changeme',
  // email = 'somebody@mail.com', first_name = 'Some', last_name = 'User');

  let user = await authenticate({username: 'my-user', password: 'changeme'});
  await screenshotAll(admin, user);
})();

/*
async function newAccount(
    {token, username, password, email, firstName, lastName}: {}) {
  let account = {
    username: username,
    password: password,
    email: email,
    first_name: first_name,
    last_name: last_name,
  } return await axios
                    .put(ADDRESS + '/api/v1/account', {account}, {
                      headers: {
                        'x-user-token': token,
                      },
                      validateStatus: false
                    })
                    .then(
                        response => {
                          if (response.status != 200 &&
                              response.data.code != ALREADY_EXISTS) {
                            console.log(response)
                            throw Error(response);
                          }
                        },
                        err => {
                          console.log(err)
                          throw Error(err);
                        });
}
*/

async function authenticate(
    {username, password}: {username: string; password: string}) {
  const response = await axios.post(
      ADDRESS + '/api/v1/login', {username, password}, {validateStatus: false});
  if (response.status != 200) {
    throw response.data;
  }
  // return {
  //   token: response.data.token, username: response.data.username,
  //       is_admin: response.data.is_admin,
  //       display_name: response.data.display_name,
  // }
}

async function screenshotAll(admin, normal) {
  const browser = await puppeteer.launch({
    bindAddress: '0.0.0.0',
    args: [
      '--headless', '--disable-gpu', '--disable-dev-shm-usage',
      '--remote-debugging-port=9222', '--remote-debugging-address=0.0.0.0'
    ]
  });

  let pages = [{url: ADDRESS + '/login', output: 'login'}];
  let users = {admin: admin, user: normal};
  for (const user in users) {
    pages.push({auth: users[user], url: ADDRESS + '/', output: 'home-' + user});
    pages.push({
      auth: users[user],
      url: ADDRESS + '/accounts/list',
      output: 'accounts-list-' + user
    });
    pages.push({
      auth: users[user],
      url: ADDRESS + '/accounts/' + users[user].username,
      output: 'accounts-edit-' + user
    });
    pages.push({
      auth: users[user],
      url: ADDRESS + '/groups/list',
      output: 'groups-list-' + user
    });
    pages.push({
      auth: users[user],
      url: ADDRESS + '/groups/users',
      output: 'groups-edit-' + user
    });
  }
  await Promise.all(pages.map(async (page) => {
    let log = 'screenshotting ' + page.url;
    if (page.auth !== undefined) log += ' for user ' + page.auth.username;
    console.log(log);
    await screenshotPage(browser, page.auth, page.url, page.output);
  }));
  console.log('done');
  await browser.close();
}

async function screenshotPage({browser, token, url, output}: {
  browser: any,
  token: Token,
  url: string,
  output: string,
}) {
  const context = await browser.createIncognitoBrowserContext();
  const page = await context.newPage();
  await page.setCacheEnabled(false);
  await page.setViewport({
    width: 860,
    height: 860,
    deviceScaleFactor: 5,
  });

  if (auth !== undefined) {
    await page.goto(url, {waitUntil: 'networkidle2'});
    await page.evaluate((auth) => {
      localStorage.setItem('x-user-token', auth.token);
      localStorage.setItem('x-user-name', auth.username);
      if (auth.is_admin) {
        localStorage.setItem('x-user-admin', 'true');
      }
      localStorage.setItem('x-user-display-name', auth.display_name);
    }, auth);
  }
  await page.goto(url, {waitUntil: 'networkidle2'});
  await page.screenshot({path: './output/' + output + '.png'});
}
