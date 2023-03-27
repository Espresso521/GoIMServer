<script>
	export let name;
	const status = fetch("/hello.json")
    .then((r) => r.json())
    .catch((err) => {
      throw new Error(err);
    });
</script>

<main>
	<h1>Hello {name}!</h1>

	<p>Visit the <a href="https://svelte.dev/tutorial">Svelte tutorial</a> to learn how to build Svelte apps.</p>

	<p class="pp">Visit the <a href="https://kotaku-blog.com">Kotaku's Blog</a> to learn how to build This website ：）</p>

	<p class="serverP">
		<strong>server respond</strong>:
		{#await status}
			loading
		{:then data}
			{data.message}
		{:catch err}
			failed to load data
		{/await}
	</p>
</main>

<style>
	main {
		text-align: center;
		padding: 1em;
		max-width: 240px;
		margin: 0 auto;
	}

	h1 {
		color: #ff3e00;
		text-transform: uppercase;
		font-size: 4em;
		font-weight: 100;
	}

	.pp {
		color: blueviolet;
	}
	
	.serverP {
		color: green;
	}

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}
</style>