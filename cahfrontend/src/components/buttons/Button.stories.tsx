import { Meta } from "@storybook/react";
import Button from "./Button";

export default {
  component: Button,
} as Meta;

export const Primary = {
  args: {
    children: "Primary",
    onClick: () => console.log("Primary button clicked"),
  },
};
