PROVIDER_REPO=intercloud
PROVIDER_PATH="$(pwd)/../"
WEBSITE_REPO=github.com/hashicorp/terraform-website

if [ ! -d "$GOPATH/src/$WEBSITE_REPO" ]
then
	echo "$WEBSITE_REPO not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone "https://$WEBSITE_REPO" "$GOPATH/src/$WEBSITE_REPO"
fi
# link the new provider repo
cd "$GOPATH/src/$WEBSITE_REPO/ext/providers"
ln -sf "$PROVIDER_PATH" "$PROVIDER_REPO"

# link the layout file
cd "$GOPATH/src/$WEBSITE_REPO/content/source/layouts"
ln -sf "../../../ext/providers/$PROVIDER_REPO/website/$PROVIDER_REPO.erb" "$PROVIDER_REPO.erb"

# link the content
cd "$GOPATH/src/$WEBSITE_REPO/content/source/docs/providers"
ln -sf "../../../../ext/providers/$PROVIDER_REPO/website/docs" "$PROVIDER_REPO"

# start middleman
cd "$PROVIDER_PATH"
make website
