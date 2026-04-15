import App from "./components/App.js";
import Home from "./components/Home.js";
import Users from "./components/Users.js";
import Layout from "./layouts/Layout.js";
import { createApp, defineRouter } from "./sugar.js";

const app = {};

createApp(document.createElement("sugar-app")).mount("#app");

const router = defineRouter([
  {
    path: "/",
    slot: "#layoutcontent",
    component: "sugar-home",
  },
  {
    path: "/users",
    slot: "#layoutcontent",
    component: "sugar-users",
  },
]);

export { app, router };
