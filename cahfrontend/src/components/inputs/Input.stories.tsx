import { Meta } from "@storybook/react";
import Input, { InputType } from "./Input";

export default {
  component: Input,
} as Meta;

export const Text = {
  args: {
    value: "",
    onchange: console.log,
    label: "Game Name",
    placeholder: "Big coconut",
    inputType: InputType.Text,
  },
};

export const Numeric = {
  args: {
    value: "",
    onchange: console.log,
    label: "Game Name",
    placeholder: "Big coconut",
    inputType: InputType.PositiveNumber,
  },
};

export const ErrorState = {
  args: {
    value: "",
    onchange: console.log,
    label: "Game Name",
    placeholder: "Big coconut",
    inputType: InputType.PositiveNumber,
    isError: true,
  },
};
