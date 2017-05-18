class UsersService:
    def __init__(self, client):
        self.client = client



    async def users_byUsername_get(self, username, headers=None, query_params=None, content_type="application/json"):
        """
        Get information on a specific user
        It is method for GET /users/{username}
        """
        uri = self.client.base_url + "/users/"+username
        return await self.client.get(uri, None, headers, query_params, content_type)


    async def users_get(self, headers=None, query_params=None, content_type="application/json"):
        """
        Get list of all developers
        It is method for GET /users
        """
        uri = self.client.base_url + "/users"
        return await self.client.get(uri, None, headers, query_params, content_type)


    async def users_post(self, data, headers=None, query_params=None, content_type="application/json"):
        """
        Add user
        It is method for POST /users
        """
        uri = self.client.base_url + "/users"
        return await self.client.post(uri, data, headers, query_params, content_type)
