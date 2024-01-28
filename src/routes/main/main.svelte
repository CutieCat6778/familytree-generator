<script lang="ts">
	import { writable } from "svelte/store";
	import {
		SvelteFlow,
		Controls,
		Background,
		BackgroundVariant,
		MiniMap,
		type SnapGrid,
	} from "@xyflow/svelte";

	import "./zyflow.style.css";
	import InputNode from "../lib/InputNode.svelte";
	import type { Data } from "../../assets/nodes";

	export let id: string;

	const nodes = writable<Data[]>();

	// same for edges
	const edges = writable([
		{ id: "a1-a2", source: "1", target: "2" },
		{ id: "a2-b", source: "2", target: "3" },
		{ id: "a2-c", source: "2", target: "4" },
	]);

	const snapGrid: SnapGrid = [25, 25];

	const nodeTypes = {
		inputNode: InputNode,
	};
</script>

<div class="h-screen">
	<SvelteFlow
		{nodes}
		{edges}
		{snapGrid}
		{nodeTypes}
		fitView
		on:nodeclick={(event) =>
			console.log("on node click", event.detail.node)}
	>
		<Controls />
		<Background variant={BackgroundVariant.Dots} />
		<MiniMap />
	</SvelteFlow>
</div>
