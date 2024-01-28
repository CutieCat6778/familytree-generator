<script lang="ts">
    import { Router, Route } from 'svelte-routing'
	import Landing from './routes/landing.svelte';
	import Main from './routes/main/main.svelte';
</script>

<Router>
    <Route component={Landing}/>
    <Route path="/:code" let:params>
        <Main id={params.id}/>
    </Route>
</Router>
