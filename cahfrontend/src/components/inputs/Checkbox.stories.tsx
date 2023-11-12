import { Meta, Story } from "@storybook/react";
import Checkbox from "./Checkbox";

export default {
  component: Checkbox,
} as Meta;

export const Checked: Story = (args) => <Checkbox {...args} />;
Checked.args = {
  label: "Example Ticked Checkbox",
  checked: true,
  onSetChecked: console.log,
};

export const UnChecked: Story = (args) => <Checkbox {...args} />;
Checked.args = {
  label: "Example Unticked Checkbox",
  checked: false,
  onSetChecked: console.log,
};
