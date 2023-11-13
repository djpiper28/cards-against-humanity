import { Meta } from "@storybook/react";
import Input from "./Input";

export default {
  component: Input,
} as Meta;

export const Checked = {
  args: {
    value: "",
    onchange: console.log,
    label: 'Game Name',
    placeholder: 'Big coconut'
  },
};
