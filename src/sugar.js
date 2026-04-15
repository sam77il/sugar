class Component extends HTMLElement {
  template;
  styles;

  constructor() {
    super();
  }

  connectedCallback() {
    this.innerHTML = this.template ?? "";
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
    get(target, p, _) {
      return target[p];
    },
  });
}

function listenToStore(eventName, cb) {
  window.addEventListener(eventName, cb);
}

function defineRouter(routes) {
  function routeHandler() {
    let foundRoute = false;
    for (const route of routes) {
      if (route.path === document.location.pathname) {
        foundRoute = true;
        const routeContent = document.querySelector("route-content");
        if (!routeContent) throw new Error("route-content component not found");
        routeContent.innerHTML = "";
        if (route.layout) {
          const layoutEl = document.createElement(route.layout);
          routeContent.appendChild(layoutEl);
          layoutEl.appendChild(document.createElement(route.component));
        } else if (route.component) {
          routeContent.appendChild(document.createElement(route.component));
        }
        break;
      }
    }
    if (!foundRoute) {
      const routeContent = document.querySelector("route-content");
      routeContent.innerHTML = "";
      routeContent.innerHTML = "404 | Oops page not found";
    }
  }
  window.addEventListener("sugarroutechange", () => {
    routeHandler();
  });

  window.addEventListener("popstate", () => {
    routeHandler();
  });

  window.addEventListener("DOMContentLoaded", () => {
    routeHandler();
  });

  const router = {
    go(path) {
      history.pushState({}, null, path);
      window.dispatchEvent(new Event("sugarroutechange"));
    },
  };
  return router;
}

class RouteContent extends Component {
  constructor() {
    super();
  }
}
defineComponent("route-content", RouteContent);

export {
  Component,
  createApp,
  defineComponent,
  defineStore,
  listenToStore,
  defineRouter,
};
