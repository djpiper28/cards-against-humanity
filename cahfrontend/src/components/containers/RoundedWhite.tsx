import { JSX } from "solid-js";

interface Props {
  children: JSX.Element;
  extraClasses?: string;
}

export default function RoundedWhite(props: Readonly<Props>) {
  return (
    <div
      class={`flex flex-col gap-5 rounded-2xl border-2 p-2 md:p-3 bg-gray-100 ${props.extraClasses}`}
    >
      {props.children}
    </div>
  );
}
