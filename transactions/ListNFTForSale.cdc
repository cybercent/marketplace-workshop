import FungibleToken from 0x01cf0e2f2f715450
import NonFungibleToken from 0x179b6b1cb6755e31
import Marketplace from 0xf3fcd2c1a78f5eee

// This transaction creates a new Sale Collection object,
// lists an NFT for sale, puts it in account storage,
// and creates a public capability to the sale so that others can buy the token.
transaction {

    prepare(acct: AuthAccount) {

        // Borrow a reference to the stored Vault
        let receiver = acct.borrow<&{FungibleToken.Receiver}>(from: /storage/MainVault)
            ?? panic("Could not borrow owner's vault reference")

        // Create a new Sale object,
        // initializing it with the reference to the owner's vault
        let sale <- Marketplace.createSaleCollection(ownerVault: receiver)

        // borrow a reference to the NFTCollection in storage
        let collectionRef = acct.borrow<&NonFungibleToken.Collection>(from: /storage/NFTCollection)
            ?? panic("Could not borrow owner's nft collection reference")

        // Withdraw the NFT from the collection that you want to sell
        // and move it into the transaction's context
        let token <- collectionRef.withdraw(withdrawID: 1)

        // List the token for sale by moving it into the sale object
        sale.listForSale(token: <-token, price: UFix64(10))

        // Store the sale object in the account storage
        acct.save(<-sale, to: /storage/NFTSale)

        // Create a public capability to the sale so that others can call its methods
        acct.link<&Marketplace.SaleCollection{Marketplace.SalePublic}>(/public/NFTSale, target: /storage/NFTSale)

        log("Sale Created for the seller's account. Selling NFT 1 for 10 tokens")
    }
}
