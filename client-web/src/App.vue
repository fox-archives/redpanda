<template>
	<div>
		<nav class="pure-g">
			<div class="pure-u-1-4 left">
				<div class="pure-menu">
					<div class="pure-menu-heading first-header">
						<img class="red-panda" src="/redpanda-icon.png" />
						<p>Red Panda</p>
					</div>
					<div class="pure-menu-list">
						<router-link to="/" class="pure-menu-link">Home</router-link>
						<router-link to="/about" class="pure-menu-link">About</router-link>
					</div>
				</div>
			</div>
			<div class="pure-u-1-24"></div>
			<div class="pure-u-17-24">
				<RouterView />
			</div>
		</nav>
	</div>
</template>

<style>
@import '@/assets/base.css';
@import 'purecss';

#app {
	max-width: 1280px;
	margin: 0 auto;
	padding: 2rem;

	font-weight: normal;
}
</style>

<style scoped>
.first-header {
	display: flex;
	align-items: center;
	gap: 10px;
}

.red-panda {
	width: 30px;
}
</style>
