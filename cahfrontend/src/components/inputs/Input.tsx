interface Props {
  value: string;
  onChanged: (text: string) => void;
  label: string;
  placeHolder?: string;
}

export default function Input(props: Readonly<Props>) {
  return (
    <label class="flex flex-row gap-2 items-center bg-white w-fit rounded-xl">
      <p class="p-2 md:px-5">{props.label}</p>
      <input
        class="rounded-xl border-2 p-1 placeholder-gray-400 font-mono h-full"
        value={props.value}
        onChange={(e) => props.onChanged(e.target.value)}
        placeholder={props.placeHolder ?? ""}
      />
    </label>
  );
}
