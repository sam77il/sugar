import { defineStore } from "../sugar.js";

const counter = {
  count: 0,
};

export default defineStore(counter, "appcounterchange");
