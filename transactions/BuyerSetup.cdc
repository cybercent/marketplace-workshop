import FungibleToken from 0x01cf0e2f2f715450
import NonFungibleToken from 0x179b6b1cb6755e31

// This transaction sets up the buyer's account for the marketplace workshop
// by creating an empty Vault, publishing a Vault reference and creating an empty NFT Collection.
transaction {
  prepare(acct: AuthAccount) {
     // Create a new empty Vault object
     let vaultA <- FungibleToken.createEmptyVault()

     // Store the vault in the account storage
     acct.save<@FungibleToken.Vault>(<-vaultA, to: /storage/MainVault)

     log("Empty Vault stored")

    // Create a public Receiver capability to the Vault
    acct.link<&FungibleToken.Vault{FungibleToken.Receiver, FungibleToken.Balance}>
             (/public/MainReceiver, target: /storage/MainVault)

    log("Created Vault references")

    // store an empty NFT Collection in account storage
    acct.save<@NonFungibleToken.Collection>(<-NonFungibleToken.createEmptyCollection(), to: /storage/NFTCollection)

    // publish a capability to the Collection in storage
    acct.link<&{NonFungibleToken.NFTReceiver}>(/public/NFTReceiver, target: /storage/NFTCollection)

    log("Created a new empty collection and published a reference")
  }
}
