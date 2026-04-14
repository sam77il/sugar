import { Component, defineComponent } from "../sugar.js";

class Users extends Component {
  template = `
  <main>
    <h1>Users Page</h1>
  </main>
  `;

  constructor() {
    super();
  }

  mounted() {
    console.log("mounted users");
  }
}

defineComponent("sugar-users", Users);

export default Users;
