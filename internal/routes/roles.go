package routes

// SetupRolesRoutes sets up the role-related routes
// func SetupRolesRoutes(router *gin.RouterGroup, rolesHandler *handler.RolesHandler) {
// 	// Group for /roles routes
// 	// rolesGroup := router.Group("/roles")
// 	// {
// 	// Public routes (if any)
// 	// }

// 	// Protected routes (require authentication)
// 	protected := router.Group("")
// 	// protected.Use(middleware.AuthMiddleware()) // Uncomment when you have auth middleware
// 	{
// 		rolesProtected := protected.Group("/roles")
// 		{
// 			rolesProtected.GET("", rolesHandler.ListRoles)
// 			rolesProtected.POST("", rolesHandler.CreateRole)
// 			rolesProtected.GET(":id", rolesHandler.GetRole)
// 			rolesProtected.PUT(":id", rolesHandler.UpdateRole)
// 			rolesProtected.DELETE(":id", rolesHandler.DeleteRole)
// 		}
// 	}
// }
