hugo-serve:
	# Since the theme use .Scratch in Hugo to implement some features, it is highly recommended that you add --disableFastRender parameter to hugo server command for the live preview of the page you are editing.
	cd guigoes.com && hugo serve --disableFastRender

hugo-new-post:
	cd guigoes.com && hugo new posts/$(title).md

hugo-build:
	cd guigoes.com && hugo --minify