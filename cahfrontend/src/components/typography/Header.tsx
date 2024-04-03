interface Props {
  text: string;
}

export default function Header(props: Readonly<Props>) {
  return <h1 class="text-2xl font-bold">{props.text}</h1>;
}
