package main

import "github.com/TutorialEdge/realtime-chat-go-react/store"

func Store_Handlers() {
	Router.HandleFunc("/store/createproduct", store.CreateProductHandler)
	Router.HandleFunc("/store/recentproducts", store.RecentProductsHandler)
	Router.HandleFunc("/store/productsbycategory", store.ProductsByCategoryHandler)
	Router.HandleFunc("/store/categories", store.AllCategoriesHandler)
	Router.HandleFunc("/store/updateproduct", store.UpdateProductHandler)
	Router.HandleFunc("/store/updateproductphoto", store.UpdateProductPhotoHandler)
	Router.HandleFunc("/store/productsbyuser", store.UserProductsHandler)
	Router.HandleFunc("/store/upgradeproductphotos", store.UpgradeProductPhotosHandler)
}
