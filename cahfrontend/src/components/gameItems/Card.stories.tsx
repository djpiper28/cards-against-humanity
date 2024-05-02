import { Meta } from "storybook-solidjs";
import Card from "./Card";

export default {
  component: Card,
} as Meta;

export const WhiteCard = {
  args: {
    isWhite: true,
    packName: "Cards Against Humanity",
    cardText: "Menstrual rage.",
  },
};

export const BlackCard = {
  args: {
    isWhite: false,
    packName: "Cards Against Humanity",
    cardText:
      "Mr. and Mrs. Diaz, we called you in because we're concerned about Cynthia. Are you aware that your daughter is ___?",
  },
};
