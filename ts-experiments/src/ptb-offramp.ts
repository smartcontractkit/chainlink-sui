/**
 * TypeScript implementation of the PTB operations from receiver_test.go
 * This script performs the same CCIP offramp execution flow using the official Sui SDK
 */

import { 
  SuiClient, 
  getFullnodeUrl,
  SuiTransactionBlockResponse,
  SuiObjectRef
} from '@mysten/sui/client';
import { Transaction } from '@mysten/sui/transactions';
import { Ed25519Keypair } from '@mysten/sui/keypairs/ed25519';
import { fromB64 } from '@mysten/sui/utils';
import * as dotenv from 'dotenv';
import { getFaucetHost, requestSuiFromFaucetV2 } from '@mysten/sui/faucet';

// Load environment variables
dotenv.config();

// Configuration interface matching the Go test environment
interface OfframpConfig {
  // Package IDs (these would be provided from your deployed contracts)
  ccipPackageId: string;
  offrampPackageId: string;
  receiverPackageId: string;
  
  // Object IDs (these would be from your deployment)
  ccipObjectRef: string;
  offrampState: string;
  
  // Network configuration
  suiUrl: string;
  
  // Signer configuration
  privateKeyB64?: string;
}

// Deployment state structure matching the Go test output
interface DeploymentState {
  network_url: string;
  chain_selectors: {
    sui_chain_selector: number;
    ethereum_chain_selector: number;
  };
  ccip_package_id: string;
  offramp_package_id: string;
  dummy_receiver_package_id: string;
  mock_link_package_id: string;
  mcms_package_id: string;
  ccip_object_ref_object_id: string;
  offramp_state_object_id: string;
  offramp_owner_cap_id: string;
  ccip_objects: {
    owner_cap_object_id: string;
    fee_quoter_cap_object_id: string;
    fee_quoter_state_object_id: string;
    nonce_manager_state_object_id: string;
    nonce_manager_cap_object_id: string;
    receiver_registry_state_object_id: string;
    rmn_remote_state_object_id: string;
    token_admin_registry_state_object_id: string;
    source_transfer_cap_object_id: string;
    dest_transfer_cap_object_id: string;
  };
  dummy_receiver_objects: {
    owner_cap_object_id: string;
    ccip_receiver_state_object_id: string;
  };
  mock_link_objects: {
    coin_metadata_object_id: string;
    treasury_cap_object_id: string;
  };
  signer_address: string;
  public_keys: string[];
  signer_addresses: string[];
  deployed_at: string;
}

// Function to load deployment state from JSON file
function loadDeploymentState(): DeploymentState | null {
  try {
    const fs = require('fs');
    const path = require('path');
    const stateFile = path.join(__dirname, '..', '..', 'state.json');
    
    if (!fs.existsSync(stateFile)) {
      console.log('⚠️  state.json not found. Run the Go test first to generate deployment state.');
      return null;
    }
    
    const stateData = fs.readFileSync(stateFile, 'utf8');
    const state = JSON.parse(stateData) as DeploymentState;
    
    console.log('✅ Loaded deployment state from state.json');
    console.log('   Deployed at:', state.deployed_at);
    console.log('   Network URL:', state.network_url);
    
    return state;
  } catch (error) {
    console.error('❌ Failed to load deployment state:', error);
    return null;
  }
}

class OfframpPTBExecutor {
  private client: SuiClient;
  private keypair: Ed25519Keypair;
  private config: OfframpConfig;

  constructor(config: OfframpConfig) {
    this.config = config;
    this.client = new SuiClient({ url: config.suiUrl });
    
    // Initialize keypair from base64 private key or generate new one
    if (config.privateKeyB64) {
      this.keypair = Ed25519Keypair.fromSecretKey(fromB64(config.privateKeyB64));
    } else {
      this.keypair = new Ed25519Keypair();
      console.log('Generated new keypair. Address:', this.keypair.toSuiAddress());
    }
  }

  /**
   * Executes the complete offramp PTB flow:
   * 1. Call DummyInitExecute (init_execute)
   * 2. Create empty token pool potatoes vector
   * 3. Call FinishExecuteWithArgs (finish_execute)
   */
  async executeOfframpPTB(): Promise<SuiTransactionBlockResponse> {
    console.log('Starting offramp PTB execution...');
    
    const txb = new Transaction();
    
    // Message parameters (matching the Go test)
    const sourceChainSelector = 2;
    const messageId = new Uint8Array(32); // 32 bytes of zeros
    const data = new TextEncoder().encode("Hello World");
    const senderAddress = this.keypair.toSuiAddress();

    console.log('Sender address:', senderAddress);
    
    // Convert sender address to bytes (remove 0x prefix and convert to bytes)
    const senderBytes = Array.from(
      new Uint8Array(Buffer.from(senderAddress.slice(2), 'hex'))
    );
    
    console.log('Execution parameters:', {
      sourceChainSelector,
      messageId: Array.from(messageId),
      senderBytes,
      data: Array.from(data),
      offrampState: this.config.offrampState
    });

    // Step 1: Call DummyInitExecute (equivalent to offrampEncoder.DummyInitExecute in Go)
    const initExecuteResult = txb.moveCall({
      target: `${this.config.offrampPackageId}::offramp::dummy_init_execute`,
      arguments: [
        txb.object(this.config.offrampState),
        txb.pure.u64(sourceChainSelector),
        txb.pure.vector('u8', Array.from(messageId)),
        txb.pure.vector('u8', senderBytes),
        txb.pure.vector('u8', Array.from(data))
      ]
    });

    console.log('Added DummyInitExecute call to PTB: ', initExecuteResult);

    // Step 2: Create empty token pool potatoes vector 
    // (equivalent to ptb.MakeMoveVec in Go)
    const hotPotatoType = `${this.config.ccipPackageId}::offramp_state_helper::CompletedDestTokenTransfer`;
    console.log('Hot potato type:', hotPotatoType);
    const tokenPoolPotatoes = txb.makeMoveVec({
      type: hotPotatoType,
      elements: []
    });

    console.log('Created empty token pool potatoes vector: ', tokenPoolPotatoes);

    // Step 3: Call FinishExecuteWithArgs (equivalent to offrampEncoder.FinishExecuteWithArgs in Go)
    txb.moveCall({
      target: `${this.config.offrampPackageId}::offramp::finish_execute`,
      arguments: [
        txb.object(this.config.offrampState),
        initExecuteResult, // The hot potato from init_execute
        tokenPoolPotatoes  // Empty vector of CompletedDestTokenTransfer
      ]
    });

    console.log('Added FinishExecuteWithArgs call to PTB');

    // Set gas budget (equivalent to the Go test's gas budget)
    txb.setGasBudget(500_000_000);

    try {
      // Sign and execute the transaction
      console.log('Signing and executing PTB...');
      const result = await this.client.signAndExecuteTransaction({
        signer: this.keypair,
        transaction: txb,
        options: {
          showEffects: true,
          showEvents: true,
          showObjectChanges: true,
          showBalanceChanges: true,
        },
        requestType: 'WaitForLocalExecution'
      });

      console.log('Transaction executed successfully!');
      console.log('Transaction digest:', result.digest);
      
      // Process and display events (equivalent to the Go test's event processing)
      await this.processTransactionEvents(result);
      
      return result;
      
    } catch (error) {
      console.error('Failed to execute PTB:', error);
      throw error;
    }
  }

  /**
   * Processes transaction events to find and decode the ExecutionCompleted event
   * (equivalent to the event processing in the Go test)
   */
  private async processTransactionEvents(result: SuiTransactionBlockResponse): Promise<void> {
    if (result.effects?.status?.status === 'success') {
      console.log('Transaction executed successfully');
      
      if (result.events && result.events.length > 0) {
        console.log(`Found ${result.events.length} events in transaction`);
        
        for (let i = 0; i < result.events.length; i++) {
          const event = result.events[i];
          console.log(`Event ${i}:`, {
            type: event.type,
            sender: event.sender,
            packageId: event.packageId,
            transactionModule: event.transactionModule,
            parsedJson: event.parsedJson
          });

          // Look for data field in the event (equivalent to Go test's data processing)
          if (event.parsedJson && typeof event.parsedJson === 'object') {
            const eventData = event.parsedJson as any;
            if (eventData.data && Array.isArray(eventData.data)) {
              // Convert array of numbers to string (equivalent to Go test's byte conversion)
              const dataBytes = new Uint8Array(eventData.data);
              const decodedString = new TextDecoder().decode(dataBytes);
              console.log('Decoded event data:', decodedString);
              
              if (decodedString === "Hello World") {
                console.log('✅ Found expected "Hello World" message in event data!');
              }
            }
          }
        }
      } else {
        console.log('⚠️  No events found in successful transaction');
      }
    } else {
      console.error('❌ Transaction failed:', result.effects?.status);
    }
  }

  /**
   * Helper method to get account balance for gas funding
   */
  async getBalance(): Promise<string> {
    const balance = await this.client.getBalance({
      owner: this.keypair.toSuiAddress()
    });
    return balance.totalBalance;
  }

    /**
   * Helper method to fund account from faucet (for testnet/devnet)
   */
  async fundFromFaucet(address: string): Promise<void> {
    try {
      console.log('Requesting funds from faucet...');
      const result = await requestSuiFromFaucetV2({
        host: getFaucetHost('localnet'),
        recipient: address,
      });
      console.log('Faucet result:', result);
      
      // Wait for the transaction to be processed
      console.log('Waiting for funds to be available...');
      await new Promise(resolve => setTimeout(resolve, 3000));
      
    } catch (error) {
      console.error('Failed to fund from faucet:', error);
    }
  }
}

// Example usage and configuration
async function main() {
  console.log('🚀 Starting CCIP Offramp PTB execution...\n');

  // Try to load deployment state from JSON file first
  const deploymentState = loadDeploymentState();
  
  let config: OfframpConfig;
  
  if (deploymentState) {
    // Use deployment state from Go test
    console.log('📋 Using deployment state from Go test execution');
    config = {
      ccipPackageId: deploymentState.ccip_package_id,
      offrampPackageId: deploymentState.offramp_package_id,
      receiverPackageId: deploymentState.dummy_receiver_package_id,
      ccipObjectRef: deploymentState.ccip_object_ref_object_id,
      offrampState: deploymentState.offramp_state_object_id,
      suiUrl: deploymentState.network_url,
      privateKeyB64: process.env.PRIVATE_KEY_B64
    };
    
    console.log('📦 Contract Configuration:');
    console.log('   CCIP Package ID:', config.ccipPackageId);
    console.log('   Offramp Package ID:', config.offrampPackageId);
    console.log('   Receiver Package ID:', config.receiverPackageId);
    console.log('   CCIP Object Ref:', config.ccipObjectRef);
    console.log('   Offramp State:', config.offrampState);
    console.log('   Network URL:', config.suiUrl);
    
  } else {
    // Fallback to environment variables or defaults
    console.log('📋 Using environment variables or defaults');
    config = {
      // These would be your actual package IDs from deployment
      ccipPackageId: process.env.CCIP_PACKAGE_ID || "0x1234567890abcdef1234567890abcdef12345678",
      offrampPackageId: process.env.OFFRAMP_PACKAGE_ID || "0x1234567890abcdef1234567890abcdef12345679", 
      receiverPackageId: process.env.RECEIVER_PACKAGE_ID || "0x1234567890abcdef1234567890abcdef12345680",
      
      // These would be your actual object IDs from deployment
      ccipObjectRef: process.env.CCIP_OBJECT_REF || "0x1234567890abcdef1234567890abcdef12345681",
      offrampState: process.env.OFFRAMP_STATE || "0x1234567890abcdef1234567890abcdef12345682",
      
      // Network configuration
      suiUrl: process.env.SUI_URL || getFullnodeUrl('localnet'),
      
      // Optional: provide private key in base64 format
      privateKeyB64: process.env.PRIVATE_KEY_B64
    };
    
    console.log('⚠️  To use actual deployed contracts, run the Go test first:');
    console.log('   cd relayer/chainwriter/ptb/offramp/receiver && go test -v');
  }

  try {
    const executor = new OfframpPTBExecutor(config);
    
    console.log('Executor initialized with address:', executor['keypair'].toSuiAddress());
    
    // Check balance
    let balance = await executor.getBalance();
    console.log('Current balance:', balance);
    
    // Fund from faucet if balance is low (for testnet/devnet)
    if (parseInt(balance) < 1000000000) { // Less than 1 SUI
      console.log('Balance low, requesting funds from faucet...');
      await executor.fundFromFaucet(executor['keypair'].toSuiAddress());
      
      // Check balance again after funding
      balance = await executor.getBalance();
      console.log('Balance after funding:', balance);
    }
    
    // Execute the offramp PTB
    const result = await executor.executeOfframpPTB();
    
    console.log('🎉 PTB execution completed successfully!');
    console.log('Final transaction digest:', result.digest);
    
  } catch (error) {
    console.error('❌ Execution failed:', error);
    process.exit(1);
  }
}

// Run the main function if this file is executed directly
if (require.main === module) {
  main().catch(console.error);
}

export { OfframpPTBExecutor, OfframpConfig };
