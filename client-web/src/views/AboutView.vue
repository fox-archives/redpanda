<template>
	<div class="about">
		<h1>About</h1>
		<p>
			RedPanda is a repsitory refactoring (<a
				target="__blank"
				href="https://martinfowler.com/bliki/RefactoringMalapropism.html"
				>restructuring</a
			>) tool. Given a list of repositories and a list of modifications, Red
			Panda allows developers to work seamlessly on multiple repositories at
			once.
		</p>
		<h2>Concepts</h2>
		<dl>
			<dt>Transaction</dt>
			<dd>A restructuring of code over multiple repositories</dd>
			<dt>Modifier</dt>
			<dd>
				A way to restructure. Can be as simple as a string replacement, regex,
				or AST thing
			</dd>
		</dl>
	</div>
</template>
