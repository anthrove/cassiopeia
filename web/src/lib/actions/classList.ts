// SOURCE: https://github.com/sveltejs/svelte/issues/3105#issuecomment-1868458487
import type { Action } from "svelte/action";

export const classList: Action<Element, string | string[]> = (node, classes) => {
  const tokens = Array.isArray(classes) ? classes : [classes];
  node.classList.add(...tokens);

  return {
    destroy() {
      node.classList.remove(...tokens);
    },
  };
};

export function setBody(classList: string|string[]):void{
    let tokens = Array.isArray(classList) ? classList : [classList];
    tokens = [...tokens,...document.body.getAttribute('data-class')?.split(' ')]
    document.body.classList.value = ""
    document.body.classList.add(...tokens);
}