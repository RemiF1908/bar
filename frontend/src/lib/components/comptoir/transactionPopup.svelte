<script lang="ts">
	import {
		ItemState,
		TransactionItemState,
		type Transaction,
		type MenuItem,
		type MenuCategory,
		type TransactionItem
	} from '$lib/api';
	import { api } from '$lib/config/config';
	import { transactionsApi } from '$lib/requests/requests';
	import { formatPrice } from '$lib/utils';
	import Error from '../error.svelte';
	import Success from '../success.svelte';
	import Transactions from './transactions.svelte';
	import { createEventDispatcher } from 'svelte';
	import { searchName } from '$lib/store/store';

	export let transaction: Transaction;
	export let close: () => void;

	let newTransaction: Transaction = structuredClone(transaction);
	let success = '';
	let error = '';

	type MenuPopup = {
		items: MenuItem[] | undefined;
		categories: MenuCategory[] | undefined;
		pickedItems: TransactionItem[] | undefined;
	};
	let menuPopup: MenuPopup | undefined;

	const eventDispatcher = createEventDispatcher();

	async function cancelTransaction() {
		let res = await transactionsApi().patchTransactionId(
			newTransaction.account_id,
			newTransaction.id,
			'canceled',
			{
				withCredentials: true
			}
		);

		if (res.status != 200) {
			error = "Une erreur s'est produite";
			setTimeout(() => {
				error = '';
			}, 1500);
			return;
		}

		transaction = newTransaction;
		success = 'Commande annulée';
		setTimeout(() => {
			success = '';
			close();
		}, 1500);
	}

	async function putBackTransaction() {
		let res = await transactionsApi().patchTransactionId(
			newTransaction.account_id,
			newTransaction.id,
			'started',
			{
				withCredentials: true
			}
		);

		if (res.status != 200) {
			error = "Une erreur s'est produite";
			setTimeout(() => {
				error = '';
			}, 1500);
			return;
		}

		transaction = newTransaction;
		success = 'Commande remise en attente';
		setTimeout(() => {
			success = '';
			close();
		}, 1500);
	}

	async function finishTransaction() {
		let res = await transactionsApi().patchTransactionId(
			newTransaction.account_id,
			newTransaction.id,
			'finished',
			{
				withCredentials: true
			}
		);

		if (res.status != 200) {
			error = "Une erreur s'est produite";
			setTimeout(() => {
				error = '';
			}, 1500);
			return;
		}

		transaction = newTransaction;
		success = 'Commande terminée';
		setTimeout(() => {
			success = '';
			searchName.set('');
			close();
		}, 1500);
	}

	async function saveTransaction() {
		for (let i = 0; i < newTransaction.items.length; i++) {
			let item = newTransaction.items[i];

			// @ts-ignore
			if (item.state == transaction.items[i].state) item.state = undefined;

			let res = await transactionsApi().patchTransactionItemId(
				newTransaction.account_id,
				newTransaction.id,
				item.item_id,
				item.state,
				item.item_amount,
				item.item_already_done,
				{
					withCredentials: true
				}
			);

			if (res.status != 200) {
				error = "Une erreur s'est produite";
				setTimeout(() => {
					error = '';
				}, 1500);
				return;
			}
		}

		if (!error) {
			transaction = newTransaction;
			success = 'Changements enregistrée';
			setTimeout(() => {
				success = '';
			}, 1500);
			reloadTransaction();
		}
	}

	function reloadTransaction() {
		transactionsApi()
			.getTransactionId(transaction.account_id, transaction.id, { withCredentials: true })
			.then((res) => {
				transaction = res.data;
				newTransaction = structuredClone(transaction);
			});
	}

	function formatTimestampToReadableDate(timestamp: number): string {
        return new Date(timestamp).toLocaleTimeString('fr-FR', {
			day: '2-digit',
			month: "2-digit",
            hour: '2-digit',
            minute: '2-digit',
			second: '2-digit'
        });
    }

	let readableDate: string = formatTimestampToReadableDate(transaction.created_at*1000);
</script>

{#if success != ''}
	<Success message={success} />
{/if}

{#if error != ''}
	<Error {error} />
{/if}

<!-- Popup overlay -->
<button
	id="overlay"
	class="absolute w-full h-full top-0 left-0 bg-black bg-opacity-50 flex justify-center items-center z-10 hover:cursor-default"
	on:click={() => {
		close();
	}}
/>

<div id="popup" class="absolute w-full h-full top-0 left-0 flex justify-center items-center">
	<!-- 
        We can modify the transaction here
        We can :
            - Check an item and mark it as completed
            - Check the transaction to mark it as completed
            - Cancel (item or transaction)
            - Undo a cancel (item or transaction)
            - Lower the amount of an item
            - Close the popup
    -->
	<div class="w-2/3 bg-white rounded-xl z-20 text-black">
			<h1 class="p-3 w-full text-xl flex-1 text-center">Commande de : 
				<span class="font-semibold">{transaction.account_name}</span>
			</h1>
			<p class="font-semibold text-xl flex-1 text-right mr-4">{readableDate}</p>
			<div class="p-5 h-full pr-4 w-full">
			<div class="grid grid-cols-6 gap-2">
				{#each newTransaction.items as item, i}
					<!-- One for each item.amount -->
					<div
						class="flex flex-col justify-center text-center break-words rounded-xl {item.state ==
						TransactionItemState.TransactionItemCanceled
							? 'bg-red-200'
							: ''} {item.state == TransactionItemState.TransactionItemFinished
							? 'bg-green-200'
							: ''}"
					>
						{item.item_name ? item.item_name : 'Test'}
						{#if item.state == TransactionItemState.TransactionItemCanceled}
							: Annulé
						{/if}
						{#if item.state == TransactionItemState.TransactionItemFinished}
							: Terminé
						{/if}
						<!-- add tooltip if menu -->
						{#if item.is_menu}
							<button
								class="flex text-white font-bold p-2 rounded justify-center"
								on:click={() => {
									menuPopup = {
										items: item.menu_items,
										categories: item.menu_categories,
										pickedItems: item.picked_categories_items
									};
								}}
							>
								<img
									src={api() + item.picture_uri}
									alt="ca charge"
									class="self-center w-10 h-10 rounded-2xl"
								/>
							</button>
						{:else}
							<img
								src={api() + item.picture_uri}
								alt="ca charge"
								class="self-center w-10 h-10 rounded-2xl"
							/>
						{/if}
						<div class="flex flex-col">
							<div>
								Total: {item.item_amount}
							</div>
							<div>
								Restant: {item.item_amount - item.item_already_done}
							</div>
						</div>
						{#if item.state == TransactionItemState.TransactionItemStarted}
							<div class="flex flex-row justify-center">
								<div class="grid grid-cols-3 gap-1">
									{#if item.item_amount > 1}
										<button
											class="bg-red-500 hover:bg-red-700 text-white font-bold p-2 rounded"
											on:click={() => {
												if (item.item_amount > 1) item.item_amount--;
											}}
										>
											<iconify-icon
												class="text-white text-lg self-center align-middle"
												icon="ic:baseline-minus"
											/>
										</button>
									{/if}
									<button
										class="bg-gray-500 hover:bg-gray-700 text-white font-bold p-2 rounded col-start-2"
										on:click={() => {
											item.state = TransactionItemState.TransactionItemCanceled;
										}}
									>
										<iconify-icon
											class="text-white text-lg self-center align-middle"
											icon="mdi:cancel"
										/>
									</button>
									<button
										class="bg-green-500 hover:bg-green-700 text-white font-bold p-2 rounded col-start-2"
										on:click={() => {
											if (item.item_already_done < item.item_amount) item.item_already_done += 1;
											if (item.item_already_done == item.item_amount)
												item.state = TransactionItemState.TransactionItemFinished;
										}}
									>
										<iconify-icon
											class="text-white text-lg self-center align-middle"
											icon="mdi:check"
										/>
									</button>
									{#if item.item_amount < transaction.items[i].item_amount}
										<button
											class="bg-green-500 hover:bg-green-700 text-white font-bold p-2 rounded col-start-3"
											on:click={() => {
												if (item.item_amount < transaction.items[i].item_amount) item.item_amount++;
											}}
										>
											<iconify-icon
												class="text-white text-lg self-center align-middle"
												icon="ic:baseline-plus"
											/>
										</button>
									{/if}
								</div>
							</div>
						{/if}
					</div>
				{/each}
			</div>
		</div>

		{#if menuPopup != undefined}
			<div class="w-full h-1 bg-gray-200" />

			<div class="w-full text-center text-lg font-bold">Ce menu contient:</div>

			<div class="grid grid-cols-5 max-h-64 p-5 gap-5">
				{#each menuPopup.categories ?? [] as cat}
					<div class="flex flex-col justify-center text-center break-words rounded-xl bg-slate-200">
						{cat.name ? cat.name : 'Test'}
						<img
							src={api() + cat.picture_uri}
							alt="..."
							class="self-center w-10 h-10 rounded-2xl"
						/>
					</div>
				{/each}
				{#each menuPopup.items ?? [] as item}
					<div class="flex flex-col justify-center text-center break-words rounded-xl bg-slate-200">
						{item.name ? item.name : 'Test'}
						<img
							src={api() + item.picture_uri}
							alt="..."
							class="self-center w-10 h-10 rounded-2xl"
						/>
					</div>
				{/each}
				{#each menuPopup.pickedItems ?? [] as pickedItem}
					<div class="flex flex-col justify-center text-center break-words rounded-xl bg-slate-200">
						{pickedItem.item_name ? pickedItem.item_name : 'Test'}
						<img
							src={api() + pickedItem.picture_uri}
							alt="..."
							class="self-center w-10 h-10 rounded-2xl"
						/>
					</div>
				{/each}
			</div>

			<div class="w-full h-1 bg-gray-200" />
		{/if}

		<div class="grid grid-cols-2 gap-4 p-8">
			<!-- Save & cancel -->
			<div class="flex flex-col gap-4 p-8">
				<button
					class="bg-green-500 rounded-xl text-white text-md font-bold p-2 h-20 w-full self-center"
					on:click={() => {
						saveTransaction();
					}}
				>
					Enregistrer les changements
				</button>
				<button
					class="bg-red-500 rounded-xl text-white text-md font-bold p-2 h-20 w-full self-center"
					on:click={() => {
						close();
					}}
				>
					Annuler les changements
				</button>
			</div>
			<div class="flex flex-col gap-4 p-8">
				<button
					class="bg-green-500 rounded-xl text-white text-md font-bold p-2 h-20 w-full"
					on:click={() => {
						finishTransaction();
					}}
				>
					Terminer la commande (paiement)
				</button>
				<button
					class="bg-gray-500 rounded-xl text-white text-md font-bold p-2 h-20 w-full"
					on:click={() => {
						cancelTransaction();
					}}
				>
					Annuler la commande (remboursement)
				</button>
				<button
					class="bg-red-500 rounded-xl text-white text-md font-bold p-2 h-20 w-full"
					on:click={() => {
						putBackTransaction();
					}}
				>
					Remettre la commande
				</button>
			</div>
		</div>

		<div class="p-5 pl-4 w-full text-lg self-center text-center">
			Prix: {formatPrice(transaction.total_cost)}
		</div>
	</div>
</div>
