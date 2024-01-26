export enum InputType {
  Text,
  PositiveNumber,
}

interface Props {
  value: string;
  onChanged: (text: string) => void;
  label: string;
  placeHolder?: string;
  inputType: InputType;
  errorState?: boolean;
  autocomplete?: string;
}

export default function Input(props: Readonly<Props>) {
  return (
    <label class="flex flex-row gap-2 items-center bg-white w-fit rounded-xl">
      <p class="p-2 md:px-5">{props.label}</p>
      <input
        id={props.label.replaceAll(" ", "-").toLowerCase()}
        class={`rounded-xl border-2 ${
          props.errorState && "border-red-600 bg-red-100"
        } p-1 placeholder-gray-400 font-mono h-full`}
        value={props.value}
        onChange={(e) => {
          let text = e.target.value;
          if (props.inputType == InputType.PositiveNumber) {
            text = text.replace(/[^\d]/, "");
          }
          props.onChanged(text);
        }}
        placeholder={props.placeHolder ?? ""}
        type={props.inputType === InputType.Text ? "text" : "number"}
        autocomplete={props.autocomplete ?? "off"}
      />
    </label>
  );
}
