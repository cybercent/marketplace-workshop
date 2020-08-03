import FungibleToken from 0x01cf0e2f2f715450
import NonFungibleToken from 0x179b6b1cb6755e31
import Marketplace from 0xf3fcd2c1a78f5eee

// This script checks  the Vault balances and NFT collections for both account
pub fun main() {
  // Get the accounts' public account objects
  let seller = getAccount(0x045a1763c93006ca)
  let buyer = getAccount(0xe03daebed8ca0615)

  // Get references to the account's receivers
  // by getting their public capability
  // and borrowing a reference from the capability
  let acct1ReceiverRef = seller.getCapability(/public/MainReceiver)!
                          .borrow<&FungibleToken.Vault{FungibleToken.Balance}>()
                                            ?? panic("Could not borrow reference to acct1 vault")
  let acct2ReceiverRef = buyer.getCapability(/public/MainReceiver)!
                          .borrow<&FungibleToken.Vault{FungibleToken.Balance}>()
                                            ?? panic("Could not borrow reference to acct2 vault")

  // Log the Vault balance of both accounts
  log("Seller's Balance")
  log(acct1ReceiverRef.balance)
  log("Buyer's Balance")
  log(acct2ReceiverRef.balance)


  // Find the public Receiver capability for their Collections
  let acct1Capability = seller.getCapability(/public/NFTReceiver)!
  let acct2Capability = buyer.getCapability(/public/NFTReceiver)!

  // borrow references from the capabilities
  let nft1Ref = acct1Capability.borrow<&{NonFungibleToken.NFTReceiver}>()
    ?? panic("Could not borrow reference to acct1 nft collection")
  let nft2Ref = acct2Capability.borrow<&{NonFungibleToken.NFTReceiver}>()
    ?? panic("Could not borrow reference to acct2 nft collection")

  // Print both collections as arrays of IDs
  log("Seller's NFTs")
  log(nft1Ref.getIDs())

  log("Buyer's NFTs")
  log(nft2Ref.getIDs())

}
