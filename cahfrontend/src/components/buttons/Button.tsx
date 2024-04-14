import { JSXElement } from "solid-js";

interface Props {
  onClick: () => void;
  children: JSXElement;
}

export default function Button(props: Props) {
  return (
    <button
      onClick={props.onClick}
      class="bg-white border-2 border-r-gray-200 rounded-2xl px-4 hover:border-yellow-500 hover:bg-gray-100"
    >
      {props.children}
    </button>
  );
}
