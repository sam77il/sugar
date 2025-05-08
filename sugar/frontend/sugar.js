function SugarLibraryInit(e) {
  let elements = document.querySelectorAll("*");
  for (let element of elements) {
    if (element.hasAttribute("sugar-onmount")) {
      let val = element.getAttribute("sugar-onmount");
      SugarOnMount(element, window[val]);
    }
    if (element.hasAttribute("sugar-click")) {
      let val = element.getAttribute("sugar-click");
      element.addEventListener("click", (e) => {
        window[val](e);
        if (element.hasAttribute("sugar-get")) {
          let val = element.getAttribute("sugar-get");
          let cb = element.getAttribute("sugar-get-cb");
          SugarGetRequest(val, window[cb]);
        }
      });
    }
    if (element.hasAttribute("sugar-onunmount")) {
      let val = element.getAttribute("sugar-onunmount");
      SugarOnUnount(element, window[val]);
    }
    if (element.hasAttribute("sugar-form-submit")) {
      let val = element.getAttribute("sugar-form-submit");

      element.addEventListener("submit", async (e) => {
        let formData = new FormData(e.target);
        let res = null;
        if (element.hasAttribute("sugar-form-pd")) {
          e.preventDefault();
          let form = e.target;
          let formData = new FormData(form);
          let params = new URLSearchParams();
          for (let [key, value] of formData.entries()) {
            params.append(key, value);
          }

          res = await fetch(form.action, {
            method: form.method,
            headers: {
              "Content-Type": "application/x-www-form-urlencoded",
            },
            body: params.toString(),
          });
        }

        window[val](res, formData);
      });
    }
  }
}

function SugarChangeState(stateName, cb) {
  let elements = document.querySelectorAll(`[sugar-state="${stateName}"]`);
  elements.forEach((element) => {
    const newContent = cb(element.innerHTML);
    element.innerHTML = newContent;
  });
}

function SugarOnMount(e, cb) {
  if (e) {
    cb(e);
  }
}

function SugarOnUnount(e, cb) {
  const observer = new MutationObserver(() => {
    if (!document.body.contains(e)) {
      cb(e);
      observer.disconnect();
    }
  });

  observer.observe(document.body, { childList: true, subtree: true });
}

function SugarCreateElement(element) {
  if (element.hasAttribute("sugar-onmount")) {
    let val = element.getAttribute("sugar-onmount");
    SugarOnMount(element, window[val]);
  }
  if (element.hasAttribute("sugar-click")) {
    let val = element.getAttribute("sugar-click");
    element.addEventListener("click", window[val]);
  }
  if (element.hasAttribute("sugar-onunmount")) {
    let val = element.getAttribute("sugar-onunmount");
    SugarOnUnount(element, window[val]);
  }
  return element;
}

async function SugarGetRequest(path, cb) {
  const res = await fetch(path, {
    method: "GET",
  });
  cb(res);
}

function abc(e) {
  SugarChangeState("test", (curVal) => {
    let val = Number(curVal);
    return val + 1;
  });
}

function buttonInit(e) {
  e.innerHTML = "Yarrak";
}

function buttonRemoved(e) {
  console.log("unmounted");
  let btn = document.createElement("button");
  btn.innerHTML = e.innerHTML;
  btn.style.backgroundColor = "red";
  btn.setAttribute("sugar-click", "abc");
  btn.setAttribute("sugar-onunmount", "buttonRemoved");
  let newBtn = SugarCreateElement(btn);
  document.body.append(newBtn);
}

async function aloalo(res) {
  let text = await res.text();
  console.log(text);
}

async function submitForm(res, formData) {
  console.log(await res.json());
  console.log(formData.get("name"));
}

document.addEventListener("DOMContentLoaded", SugarLibraryInit);
