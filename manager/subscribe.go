/*
 * Copyright 2018 The OpenWallet Authors
 * This file is part of the OpenWallet library.
 *
 * The OpenWallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The OpenWallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package manager

import (
	"github.com/blocktree/OpenWallet/openwallet"
)

//blockScanNotify 区块扫描结果通知
func (wm *WalletManager) BlockScanNotify(header *openwallet.BlockHeader) error {

	if header.Fork {
		//TODO:分叉的区块，删除提出记录

		//加载已存在所有app
		appIDs, err := wm.loadAllAppIDs()
		if err != nil {
			return err
		}

		for _, appID := range appIDs {

			wrapper, err := wm.newWalletWrapper(appID, "")
			if err != nil {
				return err
			}

			txWrapper := openwallet.NewTransactionWrapper(wrapper)
			err = txWrapper.DeleteBlockDataByHeight(header.Height)
			if err != nil {
				return err
			}

		}
		return nil
	}

	//推送数据
	for o, _ := range wm.observers {
		o.BlockScanNotify(header)
	}

	//TODO:定时删除过时的记录，保证数据库不会无限增加

	return nil
}

//BlockExtractDataNotify 区块提取结果通知
func (wm *WalletManager) BlockExtractDataNotify(sourceKey string, data *openwallet.TxExtractData) error {
	//保存提取出来的数据

	appID, accountID := wm.decodeSourceKey(sourceKey)

	wrapper, err := wm.newWalletWrapper(appID, "")
	if err != nil {
		return err
	}

	txWrapper := openwallet.NewTransactionWrapper(wrapper)
	err = txWrapper.SaveBlockExtractData(accountID, data)
	if err != nil {
		return err
	}

	//更新账户余额
	wm.RefreshAssetsAccountBalance(appID, accountID)

	account, err := wrapper.GetAssetsAccountInfo(accountID)
	if err != nil {
		return err
	}

	for o, _ := range wm.observers {
		o.BlockTxExtractDataNotify(account, data)
	}

	return nil
}

//DeleteRechargesByHeight 删除某区块高度的充值记录
func (wm *WalletManager) DeleteRechargesByHeight(height uint64) error {

	//加载已存在所有app
	appIDs, err := wm.loadAllAppIDs()
	if err != nil {
		return err
	}

	for _, appID := range appIDs {

		wrapper, err := wm.newWalletWrapper(appID, "")
		if err != nil {
			return err
		}

		txWrapper := openwallet.NewTransactionWrapper(wrapper)
		err = txWrapper.DeleteBlockDataByHeight(height)
		if err != nil {
			return err
		}

	}

	return nil
}