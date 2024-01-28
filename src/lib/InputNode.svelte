<script lang="ts">
	import { Handle, Position, type NodeProps } from "@xyflow/svelte";
	import type { Person } from "../assets/person";

	type $$Props = NodeProps;

	export let id: $$Props["id"];
	id;
	export let data: $$Props["data"] & Person;
	data;
	export let dragHandle: $$Props["dragHandle"] = undefined;
	dragHandle;
	export let type: $$Props["type"] = undefined;
	type;
	export let selected: $$Props["selected"] = undefined;
	selected;
	export let isConnectable: $$Props["isConnectable"] = undefined;
	isConnectable;
	export let zIndex: $$Props["zIndex"] = undefined;
	zIndex;
	export let width: $$Props["width"] = undefined;
	width;
	export let height: $$Props["height"] = undefined;
	height;
	export let dragging: $$Props["dragging"];
	dragging;
	export let targetPosition: $$Props["targetPosition"] = undefined;
	targetPosition;
	export let sourcePosition: $$Props["sourcePosition"] = undefined;
	sourcePosition;

    const { firstName, lastName, age, country, death, isChild, sex }: Person = data;
</script>

<div class="border-2 px-4">
    {#if isChild}
	    <Handle type="target" position={Position.Top} />
    {/if}

	<div class="">
        <span>{ lastName },</span>
        <span class="font-semibold">{ firstName }</span>
    </div>

    <div class="text-sm">
        Age: <span class="font-semibold">{age}</span><br/>
        Country: <span class="font-semibold"> {country} </span><br/>
        Sex: <span class="font-semibold"> {sex} </span>
    </div>

	<Handle type="source" position={Position.Bottom} />
</div>
