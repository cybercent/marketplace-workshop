import NonFungibleToken from 0x179b6b1cb6755e31

// This transaction mints an NFT with id=1 that is deposited into
// the sellers  NFT collection
transaction {

  // Private reference to this account's minter resource
  let minterRef: &NonFungibleToken.NFTMinter

  prepare(acct: AuthAccount) {
    // Borrow a reference for the NFTMinter in storage
    self.minterRef = acct.borrow<&NonFungibleToken.NFTMinter>(from: /storage/NFTMinter)
        ?? panic("Could not borrow owner's vault minter reference")
  }
  execute {
    // Get the seller's public account object
    let recipient = getAccount(0x045a1763c93006ca)

    // Get the Collection reference for the receiver
    // getting the public capability and borrowing a reference from it
    let receiverRef = recipient.getCapability(/public/NFTReceiver)!
                               .borrow<&{NonFungibleToken.NFTReceiver}>()
                               ?? panic("Could not borrow nft receiver reference")

    // Mint an NFT and deposit it into seller's collection
    self.minterRef.mintNFT(recipient: receiverRef)

    log("New NFT minted and deposited into seller's account")
  }
}

