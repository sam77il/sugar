import { Component, defineComponent } from "../sugar.js";

class Home extends Component {
  template = `
  <main>
    <h1>Home Page</h1>
    <sugar-users data="lolxdfxd"></sugar-users>
  </main>
  `;

  constructor() {
    super();
  }

  mounted() {
    console.log("mounted home");
  }
}

defineComponent("sugar-home", Home);

export default Home;
