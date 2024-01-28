import type { Node } from "@xyflow/svelte";
import type { Person } from "./person";

export type Data = Node & {
    id: string;
    type: "inputNode" | "group";
    data?: Person;
    position: { x: number; y: number };
}