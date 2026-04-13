class Component extends HTMLElement {
  template;
  styles;

  constructor() {
    super();
  }

  connectedCallback() {
    this.innerHTML = this.template;
    const styles = document.createElement("style");
    styles.innerHTML = this.styles ?? "";
    this.appendChild(styles);
    this.mounted();
  }

  disconnectedCallback() {
    this.unmounted();
  }
  attributeChangedCallback(name, oldValue, newValue) {
    this.attrChanged(name, oldValue, newValue);
  }

  mounted() {}
  unmounted() {}
  attrChanged(name, oldValue, newValue) {}
}

function createApp(component) {
  return {
    mount(elIdentifier) {
      const element = document.querySelector(elIdentifier);
      element.appendChild(component);
    },
  };
}

function defineComponent(tag, className) {
  customElements.define(tag, className);
}

function defineStore(obj, eventName) {
  return new Proxy(obj, {
    set(t, p, v) {
      t[p] = v;
      window.dispatchEvent(new Event(eventName));
      return true;
    },
  });
}

function listenToStore(eventName, cb) {
  window.addEventListener(eventName, cb);
}

function defineRouter(routes) {
  window.addEventListener("sugarroutechange", () => {
    let foundRoute = false;
    for (const route of routes) {
      if (route.path === document.location.pathname) {
        foundRoute = true;
        const slot = document.querySelector(route.slot);
        slot.innerHTML = "";
        slot.appendChild(document.createElement(route.component));
        break;
      }
    }
    if (!foundRoute) {
      const appEl = document.querySelector("#app");
      appEl.innerHTML = "";
      appEl.innerHTML = "404 | Oops page not found";
    }
  });

  const router = {
    go(path) {
      history.pushState({}, null, path);
      window.dispatchEvent(new Event("sugarroutechange"));
    },
  };
  return router;
}

export {
  Component,
  createApp,
  defineComponent,
  defineStore,
  listenToStore,
  defineRouter,
};
