interface Props {
  checked: boolean;
  onSetChecked: (checked: boolean) => void;
  label: string;
}

export default function Checkbox({
  checked,
  onSetChecked,
  label,
}: Readonly<Props>) {
  const background = checked ? "bg-blue-300" : "bg-white";
  return (
    <label
      class={`flex flex-row gap-3 rounded-xl border-2 w-fit p-1 px-2 ${background}`}
    >
      <input
        id={`${label} - input - checkbox`}
        type="checkbox"
        checked={checked}
        onChange={() => onSetChecked(!checked)}
        value={label}
      />
      {label}
    </label>
  );
}
