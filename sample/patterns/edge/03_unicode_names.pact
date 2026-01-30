// Edge case: Component names with numbers and underscores
// (testing identifier edge cases)

component User_Controller_V2 {
    type User_Data_V1 {
        id_number: Int
        user_name: String
    }

    provides UserAPI_V2 {
        CreateUser_V2(name: String) -> User_Data_V1
    }
}
