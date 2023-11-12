import { Meta, Story } from "@storybook/react";
import Checkbox from "./Checkbox";

export default {
  component: Checkbox,
} as Meta;

export const Checked = {
  args: {
    label: "Example Ticked Checkbox",
    checked: true,
    onSetChecked: console.log,
  },
};

export const UnChecked = {
  args: {
    label: "Example Unticked Checkbox",
    checked: false,
    onSetChecked: console.log,
  },
};
