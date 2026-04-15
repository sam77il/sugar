import { router } from "../app.js";
import countStore from "../stores/countStore.js";
import { Component, defineComponent, listenToStore } from "../sugar.js";

class App extends Component {
  template = `
  <main>
    <h1>${countStore.count}</h1>
    <button>Click Me</button>
    <a href="/users">Users</a>
    <a href="/home">Home</a>
    <a href="/lol">Test</a>
    <route-content />
  </main>
  `;
  styles = `
  main {
    background-color: #111;
  }
  main h1 {
    color: white;
  }
  `;

  constructor() {
    super();
  }

  mounted() {
    console.log("mounted app");
    listenToStore("appcounterchange", () => {
      this.querySelector("h1").innerText = countStore.count;
    });
    const btn = this.querySelector("button");
    btn.addEventListener("click", () => {
      countStore.count++;
    });

    const as = this.querySelectorAll("a");
    as.forEach((el) => {
      el.addEventListener("click", (e) => {
        e.preventDefault();
        router.go(el.href);
      });
    });
  }
}

defineComponent("sugar-app", App);

export default App;
