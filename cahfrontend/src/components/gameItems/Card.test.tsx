import { describe, it, expect } from "vitest";
import { render, screen } from "@solidjs/testing-library";
import Card, { Props } from "./Card";

describe("Card", () => {
  it("Should render the card data", () => {
    const card: Props = {
      isWhite: true,
      cardText:
        "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.",
      packName: "Poo poo",
      id: 123,
    };
    render(<Card {...card} />);
    expect(screen.getByText(card.cardText)).toBeDefined();
    expect(screen.getByText(card.packName)).toBeDefined();
    expect(screen.getByTestId("white")).toBeDefined();
  });
});
