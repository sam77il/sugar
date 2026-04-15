import { Component, defineComponent } from "../sugar.js";

class Layout extends Component {
  template = `<main id="layoutcontent"></main>`;
  constructor() {
    super();
  }
}

defineComponent("sugar-layout", Layout);

export default Layout;
