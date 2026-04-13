import { router } from "../app.js";
import countStore from "../stores/countStore.js";
import { Component, defineComponent, listenToStore } from "../sugar.js";

class App extends Component {
  template = `
  <main>
    <h1>Main Component</h1>
    <button>Click Me</button>
    <a href="/users">Users</a>
    <a href="/home">Home</a>
    <a href="/lol">Test</a>
    <div id="layoutlol"></div>
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
      this.querySelector("h1").textContent = countStore.count;
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
