interface Props {
  onSetChecked: (checked: boolean) => void;
  label: string;
  checked: boolean;
}

export default function Checkbox(props: Readonly<Props>) {
  return (
    <label
      class={`flex flex-row gap-3 rounded-xl border-2 w-fit p-1 px-2 ${
        props.checked ? "bg-blue-300" : "bg-white"
      }`}
    >
      <input
        id={`${props.label}-input-checkbox`}
        type="checkbox"
        checked={props.checked}
        onChange={() => {
          props.onSetChecked(!props.checked);
        }}
        value={props.label}
      />
      {props.label}
    </label>
  );
}
