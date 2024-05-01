interface Props {
  text: string;
}

export default function SubHeader(props: Readonly<Props>) {
  return <h2 class="text-xl font-bold">{props.text}</h2>;
}
